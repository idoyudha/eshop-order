package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/config"
	"github.com/idoyudha/eshop-order/internal/constant"
	"github.com/idoyudha/eshop-order/internal/dto"
	"github.com/idoyudha/eshop-order/internal/entity"
	"github.com/idoyudha/eshop-order/pkg/kafka"
)

type OrderCommandUseCase struct {
	repoPostgresCommand OrderPostgreCommandRepo
	producer            *kafka.ProducerServer
	warehouseService    config.WarehouseService
}

func NewOrderCommandUseCase(
	repoPostgresCommand OrderPostgreCommandRepo,
	producer *kafka.ProducerServer,
	warehouseService config.WarehouseService,
) *OrderCommandUseCase {
	return &OrderCommandUseCase{
		repoPostgresCommand,
		producer,
		warehouseService,
	}
}

type stockMovementRequest struct {
	Items   []orderItemRequest `json:"items"`
	ZipCode string             `json:"zipcode"`
}

type orderItemRequest struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int64     `json:"quantity"`
	Price       float64   `json:"price"`
}

func (u *OrderCommandUseCase) CreateOrder(ctx context.Context, order *entity.Order, token string) error {
	order.SetStatusToPending()
	err := order.GenerateOrderID()
	if err != nil {
		return fmt.Errorf("failed to generate order id: %w", err)
	}

	err = order.Address.GenerateOrderAddressID()
	if err != nil {
		return fmt.Errorf("failed to generate order address id: %w", err)
	}

	// 1. get from warehouse stock
	warehouseProductURL := fmt.Sprintf("%s/v1/stock-movements/moveout", u.warehouseService.BaseURL)
	var stockRequest stockMovementRequest
	var items []orderItemRequest

	for i := range order.Items {
		err := order.Items[i].GenerateOrderItemID()
		if err != nil {
			return fmt.Errorf("failed to generate order item id: %w", err)
		}
		items = append(items, orderItemRequest{
			ProductID: order.Items[i].ProductID,
			Quantity:  order.Items[i].ProductQuantity,
		})
	}
	stockRequest.Items = items
	stockRequest.ZipCode = order.Address.ZipCode

	requestBody, err := json.Marshal(stockRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal stock request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, warehouseProductURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create warehouse request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make warehouse request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// TODO: marshall response body to error struct
		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	return fmt.Errorf("failed to read warehouse response body: %w", err)
		// }
		// defer resp.Body.Close()

		return fmt.Errorf("warehouse service returned status: %d", resp.StatusCode)
	}

	// 2. save order to database write
	err = u.repoPostgresCommand.Insert(ctx, order)
	if err != nil {
		// TODO: handle error, send delete request to warehouse stock movement
		return fmt.Errorf("failed to insert order record: %w", err)
	}

	// 3. send event to kafka for database read
	message := dto.OrderEntityToKafkaOrderCreatedMessage(order)
	err = u.producer.Publish(
		constant.OrderCreatedTopic,
		[]byte(message.OrderID.String()),
		message,
	)
	if err != nil {
		// TODO: handle error, cancel the update if failed. or try use retry mechanism
		return fmt.Errorf("failed to produce kafka message: %w", err)
	}
	return nil
}

type kafkaProductAmountUpdated struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
}

// TODO: create kafka sale-created
// type kafkaSaleCreated struct {
// }

// func (u *OrderCommandUseCase) UpdateOrderStatus(ctx context.Context, order *entity.Order, isAcceptedPayment bool, isOrderAccepted bool) error {
// 	err := u.repoPostgresCommand.UpdateStatus(ctx, order)
// 	if err != nil {
// 		return err
// 	}

// 	if order.Status == entity.ORDER_ON_DELIVERY {
// 		if isOrderAccepted {
// 			order.SetStatusToDelivered()
// 			// TODO: send publisher to kafka sale-created
// 		}
// 		for _, item := range order.Items {
// 			message := kafkaProductAmountUpdated{
// 				ProductID: item.ProductID,
// 				Quantity:  item.ProductQuantity,
// 			}

// 			err = u.producer.Publish(
// 				constant.ProductQuantityUpdatedTopic,
// 				[]byte(item.ProductID.String()),
// 				message,
// 			)
// 			if err != nil {
// 				// TODO: handle error, cancel the update if failed. or try use retry mechanism
// 				return fmt.Errorf("failed to produce kafka message: %w", err)
// 			}
// 		}
// 	}
// 	if order.Status == entity.ORDER_PENDING {
// 		if isAcceptedPayment {
// 			order.SetStatusToOnDelivery()
// 		} else {
// 			order.SetStatusToRejected()
// 		}
// 	}

// 	return nil
// }

func (u *OrderCommandUseCase) UpdateOrderPaymentID(ctx context.Context, orderID uuid.UUID, paymentID uuid.UUID) error {
	return u.repoPostgresCommand.UpdatePaymentID(ctx, orderID, paymentID)
}
