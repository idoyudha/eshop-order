package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/internal/entity"
)

type OrderQueryUseCase struct {
	repoPostgresQuery OrderPostgreQueryRepo
}

func NewOrderQueryUseCase(repoPostgresQuery OrderPostgreQueryRepo) *OrderQueryUseCase {
	return &OrderQueryUseCase{
		repoPostgresQuery,
	}
}

func (u *OrderQueryUseCase) GetOrderByID(ctx context.Context, id uuid.UUID) (*entity.OrderView, error) {
	return u.repoPostgresQuery.GetByID(ctx, id)
}

func (u *OrderQueryUseCase) GetAllOrders(ctx context.Context) ([]*entity.OrderView, error) {
	return u.repoPostgresQuery.GetAll(ctx)
}

func (u *OrderQueryUseCase) GetOrderByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.OrderView, error) {
	return u.repoPostgresQuery.GetByUserID(ctx, userID)
}

func (u *OrderQueryUseCase) GetOrderByPaymentID(ctx context.Context, paymentID uuid.UUID) (*entity.OrderView, error) {
	return u.repoPostgresQuery.GetByPaymentID(ctx, paymentID)
}

func (u *OrderQueryUseCase) GetOrderByStatus(ctx context.Context, status string) ([]*entity.OrderView, error) {
	return u.repoPostgresQuery.GetByStatus(ctx, status)
}
