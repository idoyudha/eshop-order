package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	ORDER_VIEW_PENDING          = "PENDING"
	ORDER_VIEW_PAYMENT_ACCEPTED = "PAYMENT_ACCEPTED"
	ORDER_VIEW_ON_DELIVERY      = "ON_DELIVERY"
	ORDER_VIEW_REJECTED         = "REJECTED"
	ORDER_VIEW_DELIVERED        = "DELIVERED"
)

const (
	ORDER_VIEW_PAYMENT_PENDING  = "PENDING"
	ORDER_VIEW_PAYMENT_APPROVED = "APPROVED"
	ORDER_VIEW_PAYMENT_REJECTED = "REJECTED"
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
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	o.ID = id
	return nil
}

// admin accept payment and set the order status to ON_DELIVERY
func (o *OrderView) SetStatusToOnDelivery() {
	o.Status = ORDER_VIEW_ON_DELIVERY
}

// admin reject payment and set the order status to REJECTED
func (o *OrderView) SetStatusToRejected() {
	o.Status = ORDER_VIEW_REJECTED
}

// user just create the order and set the order status to PENDING
func (o *OrderView) SetStatusToPending() {
	o.Status = ORDER_VIEW_PENDING
	o.PaymentStatus = ORDER_VIEW_PAYMENT_PENDING
}

// user accept the delivery
func (o *OrderView) SetStatusToDelivered() {
	o.Status = ORDER_VIEW_DELIVERED
}
