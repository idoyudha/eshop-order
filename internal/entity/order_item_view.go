package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderItemView struct {
	ID                  uuid.UUID
	OrderViewID         uuid.UUID
	ProductID           uuid.UUID
	ProductName         string
	ProductPrice        float64
	ProductQuantity     int64
	ProductImageURL     string
	ProductDescription  string
	ProductCategoryID   uuid.UUID
	ProductCategoryName string
	ShippingCost        float64
	Note                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           time.Time
}

func (o *OrderItemView) GenerateOrderItemViewID() error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	o.ID = id
	return nil
}

func (o *OrderItemView) SetShippingCost(shippingCost float64) {
	o.ShippingCost = shippingCost
}
