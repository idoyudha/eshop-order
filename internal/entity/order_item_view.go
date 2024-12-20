package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderItemView struct {
	ID                  uuid.UUID
	OrderID             uuid.UUID
	ProductID           uuid.UUID
	ProductName         string
	ProductQuantity     int64
	ProductPrice        float64
	ProductImageUrl     string
	ProductDescription  string
	ProductCategoryID   uuid.UUID
	ProdcutCategoryName string
	Note                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           time.Time
}
