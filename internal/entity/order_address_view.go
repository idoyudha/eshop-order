package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderAddressView struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	Street    string
	City      string
	State     string
	ZipCode   string
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (o *OrderAddressView) GenerateOrderAddressViewID() error {
	orderID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	o.ID = orderID
	return nil
}
