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
	return u.repoPostgresCommand.Insert(ctx, order)
}

func (u *OrderCommandUseCase) UpdateOrderStatus(ctx context.Context, order *entity.Order) error {
	return u.repoPostgresCommand.UpdateStatus(ctx, order)
}

func (u *OrderCommandUseCase) UpdateOrderPaymentID(ctx context.Context, orderID uuid.UUID, paymentID uuid.UUID) error {
	return u.repoPostgresCommand.UpdatePaymentID(ctx, orderID, paymentID)
}
