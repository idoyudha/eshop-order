package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/internal/entity"
)

type OrderCommandUseCase struct {
	repoPostgresCommand OrderPostgreCommandRepo
}

func NewOrderCommandUseCase(repoPostgresCommand OrderPostgreCommandRepo) *OrderCommandUseCase {
	return &OrderCommandUseCase{
		repoPostgresCommand,
	}
}

func (u *OrderCommandUseCase) CreateOrder(ctx context.Context, order *entity.Order) error {
	order.SetStatusToPending()
	return u.repoPostgresCommand.Insert(ctx, order)
}

func (u *OrderCommandUseCase) UpdateOrderStatus(ctx context.Context, order *entity.Order, isAcceptedPayment bool, isOrderAccepted bool) error {
	err := u.repoPostgresCommand.UpdateStatus(ctx, order)
	if err != nil {
		return err
	}

	if order.Status == entity.ORDER_ON_DELIVERY {
		if isOrderAccepted {
			order.SetStatusToDelivered()
			// TODO: send publisher to kafka sale-delivered
		}
		// TODO: send publisher to kafka product-amount-updated
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
