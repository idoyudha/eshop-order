package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/internal/entity"
	"github.com/idoyudha/eshop-order/pkg/kafka"
)

const (
	ProductQuantityUpdatedTopic = "product-quantity-updated"
	SaleCreated                 = "sale-created"
)

type OrderCommandUseCase struct {
	repoPostgresCommand OrderPostgreCommandRepo
	producer            *kafka.ProducerServer
}

func NewOrderCommandUseCase(repoPostgresCommand OrderPostgreCommandRepo, producer *kafka.ProducerServer) *OrderCommandUseCase {
	return &OrderCommandUseCase{
		repoPostgresCommand,
		producer,
	}
}

func (u *OrderCommandUseCase) CreateOrder(ctx context.Context, order *entity.Order) error {
	order.SetStatusToPending()
	return u.repoPostgresCommand.Insert(ctx, order)
}

type kafkaProductAmountUpdated struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
}

// TODO: create kafka sale-created
type kafkaSaleCreated struct {
}

func (u *OrderCommandUseCase) UpdateOrderStatus(ctx context.Context, order *entity.Order, isAcceptedPayment bool, isOrderAccepted bool) error {
	err := u.repoPostgresCommand.UpdateStatus(ctx, order)
	if err != nil {
		return err
	}

	if order.Status == entity.ORDER_ON_DELIVERY {
		if isOrderAccepted {
			order.SetStatusToDelivered()
			// TODO: send publisher to kafka sale-created
		}
		for _, item := range order.Items {
			message := kafkaProductAmountUpdated{
				ProductID: item.ProductID,
				Quantity:  item.ProductQuantity,
			}

			err = u.producer.Publish(
				ProductQuantityUpdatedTopic,
				[]byte(item.ProductID.String()),
				message,
			)
			if err != nil {
				// TODO: handle error, cancel the update if failed. or try use retry mechanism
				return fmt.Errorf("failed to produce kafka message: %w", err)
			}
		}
	}
	if order.Status == entity.ORDER_PENDING {
		if isAcceptedPayment {
			order.SetStatusToOnDelivery()
		} else {
			order.SetStatusToRejected()
		}
	}

	return nil
}

func (u *OrderCommandUseCase) UpdateOrderPaymentID(ctx context.Context, orderID uuid.UUID, paymentID uuid.UUID) error {
	return u.repoPostgresCommand.UpdatePaymentID(ctx, orderID, paymentID)
}
