package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID              uuid.UUID
	OrderID         uuid.UUID
	ProductID       uuid.UUID
	ProductPrice    float64
	ProductQuantity int64
	Note            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}
