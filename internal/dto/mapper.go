package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/internal/entity"
)

func OrderEntityToKafkaOrderCreatedMessage(order *entity.Order) KafkaOrderCreated {
	return KafkaOrderCreated{
		OrderID:    order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		Items:      orderItemEntityToKafkaOrderItemsCreated(order.Items),
		Address: KafkaOrderAddressCreated{
			OrderID: order.Address.OrderID,
			Street:  order.Address.Street,
			City:    order.Address.City,
			State:   order.Address.State,
			ZipCode: order.Address.ZipCode,
			Note:    order.Address.Note,
		},
	}
}

func orderItemEntityToKafkaOrderItemsCreated(items []entity.OrderItem) []KafkaOrderItemsCreated {
	var kafkaItems []KafkaOrderItemsCreated
	for _, item := range items {
		kafkaItems = append(kafkaItems, KafkaOrderItemsCreated{
			OrderID:         item.OrderID,
			ProductID:       item.ProductID,
			ProductQuantity: item.ProductQuantity,
			ShippingCost:    item.ShippingCost,
			Note:            item.Note,
		})
	}

	return kafkaItems
}

func PaymentMessageUpdateToOrderEntity(msg KafkaPaymentUpdated) entity.Order {
	return entity.Order{
		ID:        msg.OrderID,
		PaymentID: msg.PaymentID,
	}
}

func PaymentMessageToOrderViewEntity(message KafkaPaymentUpdated) entity.OrderView {
	return entity.OrderView{
		OrderID:          message.OrderID,
		PaymentID:        message.PaymentID,
		PaymentStatus:    message.Status,
		PaymentImageURL:  message.ImageURL,
		PaymentAdminNote: message.Note,
		UpdatedAt:        time.Now(),
	}
}

func OrderEntityToKafkaOrderStatusUpdatedMessage(order *entity.Order) KafkaOrderStatusUpdated {
	return KafkaOrderStatusUpdated{
		OrderID: order.ID,
		Status:  order.Status,
	}
}

func OrderStatusUpdatedMessageToOrderViewEntity(msg KafkaOrderStatusUpdated) entity.OrderView {
	return entity.OrderView{
		OrderID:   msg.OrderID,
		Status:    msg.Status,
		UpdatedAt: time.Now(),
	}
}

func OrderEntityToKafkaSaleCreatedMessage(order *entity.Order, products map[uuid.UUID]float64) KafkaSaleCreated {
	return KafkaSaleCreated{
		OrderID: order.ID,
		UserID:  order.UserID,
		Items:   orderItemEntityToKafkaSaleItemsCreated(order.Items, products),
	}
}

func orderItemEntityToKafkaSaleItemsCreated(items []entity.OrderItem, products map[uuid.UUID]float64) []KafkaSaleItemCreated {
	var kafkaItems []KafkaSaleItemCreated
	for _, item := range items {
		kafkaItems = append(kafkaItems, KafkaSaleItemCreated{
			ProductID: item.ProductID,
			Quantity:  item.ProductQuantity,
			Price:     products[item.ProductID],
		})
	}

	return kafkaItems
}
