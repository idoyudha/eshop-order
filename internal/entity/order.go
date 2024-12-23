package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	ORDER_PENDING          = "PENDING"
	ORDER_PAYMENT_ACCEPTED = "PAYMENT_ACCEPTED"
	ORDER_ON_DELIVERY      = "ON_DELIVERY"
	ORDER_REJECTED         = "REJECTED"
	ORDER_DELIVERED        = "DELIVERED"
)

type Order struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Status     string
	TotalPrice float64
	PaymentID  uuid.UUID
	Items      []OrderItem
	Address    OrderAddress
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

func (o *Order) GenerateOrderID() error {
	orderID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	o.ID = orderID
	return nil
}

// admin accept payment and set the order status to ON_DELIVERY
func (o *Order) SetStatusToOnDelivery() {
	o.Status = ORDER_ON_DELIVERY
}

// admin reject payment and set the order status to REJECTED
func (o *Order) SetStatusToRejected() {
	o.Status = ORDER_REJECTED
}

// user just create the order and set the order status to PENDING
func (o *Order) SetStatusToPending() {
	o.Status = ORDER_PENDING
}

// user accept the delivery
func (o *Order) SetStatusToDelivered() {
	o.Status = ORDER_DELIVERED
}
