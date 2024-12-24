package dto

import "github.com/idoyudha/eshop-order/internal/entity"

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
			Note:            item.Note,
		})
	}

	return kafkaItems
}
