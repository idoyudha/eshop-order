package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/config"
	"github.com/idoyudha/eshop-order/internal/constant"
	"github.com/idoyudha/eshop-order/internal/dto"
	"github.com/idoyudha/eshop-order/internal/entity"
	"github.com/idoyudha/eshop-order/internal/usecase"
	kafkaConSrv "github.com/idoyudha/eshop-order/pkg/kafka"
	"github.com/idoyudha/eshop-order/pkg/logger"
)

type kafkaConsumerRoutes struct {
	ucoq usecase.OrderQuery
	l    logger.Interface
	p    config.ProductService
}

func KafkaNewRouter(
	ucoq usecase.OrderQuery,
	l logger.Interface,
	c *kafkaConSrv.ConsumerServer,
	p config.ProductService,
) error {
	routes := &kafkaConsumerRoutes{
		ucoq: ucoq,
		l:    l,
		p:    p,
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
			return nil
		default:
			// l.Debug("Attempting to read message...")
			ev, err := c.Consumer.ReadMessage(3 * time.Second)
			if err != nil {
				// log.Println("CONSUME CART SERVICE!!")
				// Errors are informational and automatically handled by the consumer
				if kerr, ok := err.(kafka.Error); ok && kerr.Code() == kafka.ErrTimedOut {
					// l.Debug("Timeout waiting for message, continuing...")
					continue
				}
				l.Error("Error reading message: ", err)
				continue
			}

			switch *ev.TopicPartition.Topic {
			case constant.OrderCreatedTopic:
				if err := routes.handleOrderViewCreated(ev); err != nil {
					l.Error("Failed to handle order view creted: %w", err)
				}
			default:
				l.Info("Unknown topic: %s", *ev.TopicPartition.Topic)
			}

			log.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}

	return nil
}

type getProductResponse struct {
	ID          string  `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CategoryID  string  `json:"category_id"`
}

type restSuccess struct {
	Code    int                `json:"code"`
	Data    getProductResponse `json:"data"`
	Message string             `json:"message"`
}

func (r *kafkaConsumerRoutes) handleOrderViewCreated(msg *kafka.Message) error {
	var message dto.KafkaOrderCreated

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleOrderViewCreated")
		return err
	}

	// get product data from product service
	var items []entity.OrderItemView
	for _, item := range message.Items {
		productServiceURL := fmt.Sprintf("%s/v1/products/%s", r.p.BaseURL, item.ProductID)
		req, err := http.NewRequest(http.MethodGet, productServiceURL, nil)
		if err != nil {
			r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleOrderViewCreated")
			return fmt.Errorf("failed to create product request: %w", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to make product request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to make product request, status not OK: %w", err)
		}

		var restSuccess restSuccess
		if err := json.NewDecoder(resp.Body).Decode(&restSuccess); err != nil {
			return fmt.Errorf("failed to decode product response: %w", err)
		}

		items = append(items, entity.OrderItemView{
			ProductName:        restSuccess.Data.Name,
			ProductImageURL:    restSuccess.Data.ImageURL,
			ProductDescription: restSuccess.Data.Description,
			ProductCategoryID:  uuid.MustParse(restSuccess.Data.CategoryID),
			ProductPrice:       restSuccess.Data.Price,
		})
	}

	orderViewEntity := kafkaOrderCreatedToOrderView(&message)
	orderViewEntity.Items = items
	err := r.ucoq.CreateOrderView(context.Background(), &orderViewEntity)
	if err != nil {
		r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleOrderViewCreated")
		return fmt.Errorf("failed to create order view: %w", err)
	}

	return nil
}

func kafkaOrderCreatedToOrderView(msg *dto.KafkaOrderCreated) entity.OrderView {
	return entity.OrderView{
		OrderID:    msg.OrderID,
		UserID:     msg.UserID,
		TotalPrice: msg.TotalPrice,
		Address: entity.OrderAddressView{
			Street:    msg.Address.Street,
			City:      msg.Address.City,
			State:     msg.Address.State,
			ZipCode:   msg.Address.ZipCode,
			Note:      msg.Address.Note,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
