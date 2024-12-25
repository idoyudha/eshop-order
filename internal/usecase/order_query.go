package usecase

import (
	"context"
	"fmt"

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

func (u *OrderQueryUseCase) CreateOrderView(ctx context.Context, order *entity.OrderView) error {
	order.SetStatusToPending()

	err := order.GenerateOrderViewID()
	if err != nil {
		return fmt.Errorf("failed to generate order view id: %w", err)
	}

	for i := range order.Items {
		err := order.Items[i].GenerateOrderItemViewID()
		if err != nil {
			return fmt.Errorf("failed to generate order view item id: %w", err)
		}
	}

	err = order.Address.GenerateOrderAddressViewID()
	if err != nil {
		return fmt.Errorf("failed to generate order view address id: %w", err)
	}

	return u.repoPostgresQuery.Insert(ctx, order)
}

func (u *OrderQueryUseCase) UpdateOrderViewPayment(ctx context.Context, order *entity.OrderView) error {
	return u.repoPostgresQuery.UpdatePayment(ctx, order)
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
