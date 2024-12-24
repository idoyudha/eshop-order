package dto

import "github.com/google/uuid"

// this dto is used to send order created event to kafka in the same service, to prevent repetitive struct
type KafkaOrderCreated struct {
	OrderID    uuid.UUID                `json:"order_id"`
	UserID     uuid.UUID                `json:"user_id"`
	TotalPrice float64                  `json:"total_price"`
	Items      []KafkaOrderItemsCreated `json:"items"`
	Address    KafkaOrderAddressCreated `json:"address"`
}

type KafkaOrderItemsCreated struct {
	OrderID         uuid.UUID `json:"order_id"`
	ProductID       uuid.UUID `json:"product_id"`
	ProductQuantity int64     `json:"product_quantity"`
	Note            string    `json:"note"`
}

type KafkaOrderAddressCreated struct {
	OrderID uuid.UUID `json:"order_id"`
	Street  string    `json:"street"`
	City    string    `json:"city"`
	State   string    `json:"state"`
	ZipCode string    `json:"zipcode"`
	Note    string    `json:"note"`
}
