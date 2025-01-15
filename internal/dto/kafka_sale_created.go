package dto

import "github.com/google/uuid"

type KafkaSaleCreated struct {
	OrderID uuid.UUID              `json:"order_id"`
	UserID  uuid.UUID              `json:"user_id"`
	Items   []KafkaSaleItemCreated `json:"items"`
}

type KafkaSaleItemCreated struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
}
