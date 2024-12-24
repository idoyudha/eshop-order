package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderView struct {
	ID               uuid.UUID
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
