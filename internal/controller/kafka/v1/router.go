package v1

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/internal/usecase"
	kafkaConSrv "github.com/idoyudha/eshop-order/pkg/kafka"
	"github.com/idoyudha/eshop-order/pkg/logger"
)

type kafkaConsumerRoutes struct {
	uc usecase.OrderCommand
	l  logger.Interface
}

func KafkaNewRouter(
	uc usecase.OrderCommand,
	l logger.Interface,
	c *kafkaConSrv.ConsumerServer,
) error {
	routes := &kafkaConsumerRoutes{uc, l}

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
			ev, err := c.Consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// log.Println("CONSUME CART SERVICE!!")
				// Errors are informational and automatically handled by the consumer
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				l.Error("Error reading message: ", err)
				continue
			}

			switch *ev.TopicPartition.Topic {
			case kafkaConSrv.PaymentCreatedTopic:
				if err := routes.handlePaymentCreated(ev); err != nil {
					l.Error("Failed to handle product creation: %w", err)
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

type kafkaPaymentCreatedMessage struct {
	OrderID   uuid.UUID `json:"order_id"`
	PaymentID uuid.UUID `json:"payment_id"`
}

func (r *kafkaConsumerRoutes) handlePaymentCreated(msg *kafka.Message) error {
	var message kafkaPaymentCreatedMessage
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleProductCreated")
		return err
	}

	if err := r.uc.UpdateOrderPaymentID(context.Background(), message.OrderID, message.PaymentID); err != nil {
		return err
	}

	return nil
}
