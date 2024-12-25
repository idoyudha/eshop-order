package dto

import "github.com/google/uuid"

type KafkaPaymentUpdated struct {
	PaymentID uuid.UUID `json:"paymentId"`
	OrderID   uuid.UUID `json:"orderId"`
	ImageURL  string    `json:"imageUrl"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
}
