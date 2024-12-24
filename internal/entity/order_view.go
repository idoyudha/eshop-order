package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderView struct {
	ID               uuid.UUID
	OrderID          uuid.UUID
	UserID           uuid.UUID
	Status           string
	TotalPrice       float64
	PaymentID        uuid.UUID
	PaymentStatus    string
	PaymentImageURL  string
	PaymentAdminNote string
	Items            []OrderItemView
	Address          OrderAddressView
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}

func (o *OrderView) GenerateOrderViewID() error {
	orderID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	o.ID = orderID
	return nil
}
