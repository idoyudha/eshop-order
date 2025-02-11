package v1

import (
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/internal/entity"
)

func CreateOrderRequestToOrderEntity(req createOrderRequest, userID uuid.UUID) (entity.Order, error) {
	orderID, err := uuid.NewV7()
	if err != nil {
		return entity.Order{}, err
	}
	var items []entity.OrderItem
	var totalPice float64 // TODO: should be get from warehouse price for accurate
	for _, item := range req.Items {
		items = append(items, entity.OrderItem{
			OrderID:         orderID,
			ProductID:       item.ProductID,
			ProductQuantity: item.Quantity,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		})
		totalPice += float64(item.Quantity) * item.Price
	}

	return entity.Order{
		ID:         orderID,
		UserID:     userID,
		TotalPrice: totalPice,
		PaymentID:  uuid.UUID{},
		Items:      items,
		Address: entity.OrderAddress{
			OrderID:   orderID,
			Street:    req.Address.Street,
			City:      req.Address.City,
			State:     req.Address.State,
			ZipCode:   req.Address.ZipCode,
			Note:      req.Address.Note,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func UpdateOrderRequestToOrderEntity(orderID uuid.UUID) entity.Order {
	return entity.Order{
		ID:        orderID,
		UpdatedAt: time.Now(),
	}
}

func OrderEntityToCreatedOrderResponse(order entity.Order) orderResponse {
	var items []itemsOrderResponse
	for _, item := range order.Items {
		items = append(items, itemsOrderResponse{
			OrderID:      item.OrderID,
			ProductID:    item.ProductID,
			Quantity:     item.ProductQuantity,
			ShippingCost: item.ShippingCost,
			Note:         item.Note,
		})
	}

	return orderResponse{
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		Items:      items,
		Address: addressOrderResponse{
			OrderID: order.Address.OrderID,
			Street:  order.Address.Street,
			City:    order.Address.City,
			State:   order.Address.State,
			ZipCode: order.Address.ZipCode,
			Note:    order.Address.Note,
		},
	}
}

func OrderViewEntityToGetManyOrderResponse(orders []*entity.OrderView) []orderResponse {
	var res []orderResponse
	for _, order := range orders {
		var items []itemsOrderResponse
		for _, item := range order.Items {
			items = append(items, itemsOrderResponse{
				OrderID:      item.OrderViewID,
				ProductID:    item.ProductID,
				ProductName:  item.ProductName,
				ImageURL:     item.ProductImageURL,
				Price:        item.ProductPrice,
				Quantity:     item.ProductQuantity,
				ShippingCost: item.ShippingCost,
				Note:         item.Note,
			})
		}

		res = append(res, orderResponse{
			ID:              order.OrderID,
			Status:          order.Status,
			TotalPrice:      order.TotalPrice,
			PaymentID:       order.PaymentID,
			PaymentStatus:   order.PaymentStatus,
			PaymentImageURL: order.PaymentImageURL,
			Items:           items,
			Address: addressOrderResponse{
				OrderID: order.Address.OrderViewID,
				Street:  order.Address.Street,
				City:    order.Address.City,
				State:   order.Address.State,
				ZipCode: order.Address.ZipCode,
				Note:    order.Address.Note,
			},
			CreatedAt: order.CreatedAt,
		})
	}
	return res
}

func OrderViewEntityToGetOneOrderResponse(order *entity.OrderView) orderResponse {
	var items []itemsOrderResponse
	for _, item := range order.Items {
		items = append(items, itemsOrderResponse{
			OrderID:      item.OrderViewID,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ImageURL:     item.ProductImageURL,
			Price:        item.ProductPrice,
			Quantity:     item.ProductQuantity,
			ShippingCost: item.ShippingCost,
			Note:         item.Note,
		})
	}

	return orderResponse{
		ID:              order.OrderID,
		Status:          order.Status,
		TotalPrice:      order.TotalPrice,
		PaymentID:       order.PaymentID,
		PaymentStatus:   order.PaymentStatus,
		PaymentImageURL: order.PaymentImageURL,
		Items:           items,
		Address: addressOrderResponse{
			OrderID: order.Address.OrderViewID,
			Street:  order.Address.Street,
			City:    order.Address.City,
			State:   order.Address.State,
			ZipCode: order.Address.ZipCode,
			Note:    order.Address.Note,
		},
		CreatedAt: order.CreatedAt,
	}
}
