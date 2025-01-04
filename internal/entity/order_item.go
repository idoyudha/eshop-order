package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID              uuid.UUID
	OrderID         uuid.UUID
	ProductID       uuid.UUID
	ProductQuantity int64
	Note            string
	ShippingCost    float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

func (o *OrderItem) GenerateOrderItemID() error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	o.ID = id
	return nil
}

func (o *OrderItem) SetShippingCost(shippingCost float64) {
	o.ShippingCost = shippingCost
}
