package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID              uuid.UUID
	OrderID         uuid.UUID
	ProductID       uuid.UUID
	ProductQuantity int64
	Note            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}
