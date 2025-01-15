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
		UpdatePaymentID(context.Context, *entity.Order) error
	}

	OrderPostgreQueryRepo interface {
		Insert(context.Context, *entity.OrderView) error
		UpdatePayment(context.Context, *entity.OrderView) error
		GetByID(context.Context, uuid.UUID) (*entity.OrderView, error)
		GetAll(context.Context) ([]*entity.OrderView, error)
		GetByUserID(context.Context, uuid.UUID) ([]*entity.OrderView, error)
		GetByPaymentID(context.Context, uuid.UUID) (*entity.OrderView, error)
		GetByStatus(context.Context, string) ([]*entity.OrderView, error)
		UpdateStatus(context.Context, *entity.OrderView) error
	}

	OrderCommand interface {
		CreateOrder(context.Context, *entity.Order, string) error
		UpdateOrderStatus(context.Context, *entity.Order, string) error
		UpdateOrderPaymentID(context.Context, *entity.Order, string) error
	}

	OrderQuery interface {
		CreateOrderView(context.Context, *entity.OrderView) error
		UpdateOrderViewPayment(context.Context, *entity.OrderView, string) error
		GetOrderByID(context.Context, uuid.UUID) (*entity.OrderView, error)
		GetAllOrders(context.Context) ([]*entity.OrderView, error)
		GetOrderByUserID(context.Context, uuid.UUID) ([]*entity.OrderView, error)
		GetOrderByPaymentID(context.Context, uuid.UUID) (*entity.OrderView, error)
		GetOrderByStatus(context.Context, string) ([]*entity.OrderView, error)
		UpdateOrderStatus(context.Context, *entity.OrderView, string) error
	}
)
