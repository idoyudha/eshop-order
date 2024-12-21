package usecase

import (
	"context"

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

func (u *OrderCommandUseCase) UpdateStatus(ctx context.Context, order *entity.Order) error {
	return u.repoPostgresCommand.UpdateStatus(ctx, order)
}
