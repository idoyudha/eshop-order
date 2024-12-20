package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/internal/entity"
)

type (
	OrderPostgreCommandRepo interface {
		Insert(context.Context, *entity.Order) error
		UpdateStatus(context.Context, *entity.Order) error
	}

	OrderPostgreQueryRepo interface {
		GetByID(context.Context, uuid.UUID) (*entity.Order, error)
		GetAll(context.Context) ([]*entity.Order, error)
		GetByUserID(context.Context, uuid.UUID) ([]*entity.Order, error)
		GetByPaymentID(context.Context, uuid.UUID) ([]*entity.Order, error)
		GetByStatus(context.Context, string) ([]*entity.Order, error)
	}

	OrderCommand interface {
		CreateOrder(context.Context, *entity.Order) error
		UpdateOrderStatus(context.Context, *entity.Order) error
	}

	OrderQuery interface {
		GetOrderByID(context.Context, uuid.UUID) (*entity.Order, error)
		GetAllOrders(context.Context) ([]*entity.Order, error)
		GetOrderByUserID(context.Context, uuid.UUID) ([]*entity.Order, error)
		GetOrderByPaymentID(context.Context, uuid.UUID) ([]*entity.Order, error)
		GetOrderByStatus(context.Context, string) ([]*entity.Order, error)
	}
)
