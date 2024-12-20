package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	ON_DELIVERY = "ON_DELIVERY"
	REJECTED    = "REJECTED"
)

type Order struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Status     string
	TotalPrice float64
	PaymentID  uuid.UUID
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
