package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	PENDING     = "PENDING"
	ON_DELIVERY = "ON_DELIVERY"
	REJECTED    = "REJECTED"
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

// admin accept payment and set the order status to ON_DELIVERY
func (o *Order) SetStatusToOnDelivery() {
	o.Status = ON_DELIVERY
}

// admin reject payment and set the order status to REJECTED
func (o *Order) SetStatusToRejected() {
	o.Status = REJECTED
}

// user just create the order and set the order status to PENDING
func (o *Order) SetStatusTopPending() {
	o.Status = PENDING
}
