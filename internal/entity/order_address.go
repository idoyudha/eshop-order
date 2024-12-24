package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderAddress struct {
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

func (o *OrderAddress) GenerateOrderAddressID() error {
	orderID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	o.ID = orderID
	return nil
}
