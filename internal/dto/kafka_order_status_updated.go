package dto

import "github.com/google/uuid"

type KafkaOrderStatusUpdated struct {
	OrderID uuid.UUID `json:"orderId"`
	Status  string    `json:"status"`
}
