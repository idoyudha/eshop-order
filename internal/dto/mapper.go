package dto

import (
	"time"

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
