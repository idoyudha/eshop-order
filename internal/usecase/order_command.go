package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/config"
	"github.com/idoyudha/eshop-order/internal/constant"
	"github.com/idoyudha/eshop-order/internal/dto"
	"github.com/idoyudha/eshop-order/internal/entity"
	"github.com/idoyudha/eshop-order/pkg/kafka"
)

type OrderCommandUseCase struct {
	repoPostgresCommand OrderPostgreCommandRepo
	producer            *kafka.ProducerServer
	warehouseService    config.WarehouseService
	shippingCostService config.ShippingCostService
}

func NewOrderCommandUseCase(
	repoPostgresCommand OrderPostgreCommandRepo,
	producer *kafka.ProducerServer,
	warehouseService config.WarehouseService,
	shippingCostService config.ShippingCostService,
) *OrderCommandUseCase {
	return &OrderCommandUseCase{
		repoPostgresCommand,
		producer,
		warehouseService,
		shippingCostService,
	}
}

type stockMovementRequest struct {
	Items   []orderItemRequest `json:"items"`
	ZipCode string             `json:"zipcode"`
}

type orderItemRequest struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int64     `json:"quantity"`
	Price       float64   `json:"price"`
}

func createStockMovement(ctx context.Context, whBaseURL string, stockRequest stockMovementRequest, token string) error {
	warehouseProductURL := fmt.Sprintf("%s/v1/stock-movements/moveout", whBaseURL)
	requestBody, err := json.Marshal(stockRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal stock request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, warehouseProductURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create warehouse request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make warehouse request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// TODO: marshall response body to error struct
		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	return fmt.Errorf("failed to read warehouse response body: %w", err)
		// }
		// defer resp.Body.Close()
		return fmt.Errorf("warehouse service returned status: %d", resp.StatusCode)
	}

	return nil
}

type getNearestWarehouseRequest struct {
	ZipCode   string    `json:"zip_code"`
	ProductID uuid.UUID `json:"product_id"`
}

type getNearestWarehouseResponse struct {
	Code int `json:"code"`
	Data struct {
		ZipCode string `json:"zip_code"`
	}
	Message string `json:"message"`
}

func getNearestWarehouse(ctx context.Context, token, whBaseURL, zipCode string, productID uuid.UUID) (*string, error) {
	var request getNearestWarehouseRequest
	request.ZipCode = zipCode
	request.ProductID = productID

	nearestWarehouseURL := fmt.Sprintf("%s/v1/warehouse-products/nearest", whBaseURL)
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal nearest warehouse request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, nearestWarehouseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create nearest warehouse request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make nearest warehouse request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nearest warehouse service returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read nearest warehouse response body: %w", err)
	}
	defer resp.Body.Close()

	var nearestWarehouseResponse getNearestWarehouseResponse
	if err := json.Unmarshal(body, &nearestWarehouseResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nearest warehouse response: %w", err)
	}

	return &nearestWarehouseResponse.Data.ZipCode, nil
}

type shippingCostRequest struct {
	FromZip string `json:"from_zip"`
	ToZip   string `json:"to_zip"`
}

type shippingCostResponse struct {
	Code int `json:"code"`
	Data struct {
		ShippingCost float64 `json:"shipping_cost"`
	}
	Message string `json:"message"`
}

func getShippingCost(ctx context.Context, scBaseURL string, request shippingCostRequest) (float64, error) {
	shippingCostURL := fmt.Sprintf("%s/shipping-cost", scBaseURL)
	requestBody, err := json.Marshal(request)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal shipping cost request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, shippingCostURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, fmt.Errorf("failed to create shipping cost request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make shipping cost request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("shipping cost service returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read shipping cost response body: %w", err)
	}
	defer resp.Body.Close()

	var shippingCostResponse shippingCostResponse
	err = json.Unmarshal(body, &shippingCostResponse)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal shipping cost response: %w", err)
	}

	return shippingCostResponse.Data.ShippingCost, nil
}

func (u *OrderCommandUseCase) CreateOrder(ctx context.Context, order *entity.Order, token string) error {
	order.SetStatusToPending()
	err := order.GenerateOrderID()
	if err != nil {
		return fmt.Errorf("failed to generate order id: %w", err)
	}

	err = order.Address.GenerateOrderAddressID()
	if err != nil {
		return fmt.Errorf("failed to generate order address id: %w", err)
	}
	// 1. create stock movement
	var stockRequest stockMovementRequest
	var items []orderItemRequest

	for i := range order.Items {
		err := order.Items[i].GenerateOrderItemID()
		if err != nil {
			return fmt.Errorf("failed to generate order item id: %w", err)
		}

		// 2. get and set shipping cost
		nearestZipCode, err := getNearestWarehouse(ctx, token, u.warehouseService.BaseURL, order.Address.ZipCode, order.Items[i].ProductID)
		if err != nil {
			return fmt.Errorf("failed to get nearest warehouse zipcode: %w", err)
		}

		shippingCostReq := shippingCostRequest{
			FromZip: order.Address.ZipCode,
			ToZip:   *nearestZipCode,
		}

		shippingCost, err := getShippingCost(ctx, u.shippingCostService.URL, shippingCostReq)
		if err != nil {
			return fmt.Errorf("failed to get shipping cost: %w", err)
		}

		order.Items[i].SetShippingCost(shippingCost)

		order.AddShippingCost(shippingCost)
		items = append(items, orderItemRequest{
			ProductID: order.Items[i].ProductID,
			Quantity:  order.Items[i].ProductQuantity,
		})
	}
	stockRequest.Items = items
	stockRequest.ZipCode = order.Address.ZipCode
	err = createStockMovement(ctx, u.warehouseService.BaseURL, stockRequest, token)
	if err != nil {
		return fmt.Errorf("failed to create stock movement: %w", err)
	}

	// 3. save order to database write
	err = u.repoPostgresCommand.Insert(ctx, order)
	if err != nil {
		// TODO: handle error, send delete request to warehouse stock movement
		return fmt.Errorf("failed to insert order record: %w", err)
	}

	// 4. send event to kafka for database read
	message := dto.OrderEntityToKafkaOrderCreatedMessage(order)
	err = u.producer.Publish(
		constant.OrderCreatedTopic,
		[]byte(message.OrderID.String()),
		message,
	)
	if err != nil {
		// TODO: handle error, cancel the update if failed. or try use retry mechanism
		return fmt.Errorf("failed to produce kafka message: %w", err)
	}
	return nil
}

type kafkaProductAmountUpdated struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
}

// TODO: create kafka sale-created
// type kafkaSaleCreated struct {
// }

// func (u *OrderCommandUseCase) UpdateOrderStatus(ctx context.Context, order *entity.Order, isAcceptedPayment bool, isOrderAccepted bool) error {
// 	err := u.repoPostgresCommand.UpdateStatus(ctx, order)
// 	if err != nil {
// 		return err
// 	}

// 	if order.Status == entity.ORDER_ON_DELIVERY {
// 		if isOrderAccepted {
// 			order.SetStatusToDelivered()
// 			// TODO: send publisher to kafka sale-created
// 		}
// 		for _, item := range order.Items {
// 			message := kafkaProductAmountUpdated{
// 				ProductID: item.ProductID,
// 				Quantity:  item.ProductQuantity,
// 			}

// 			err = u.producer.Publish(
// 				constant.ProductQuantityUpdatedTopic,
// 				[]byte(item.ProductID.String()),
// 				message,
// 			)
// 			if err != nil {
// 				// TODO: handle error, cancel the update if failed. or try use retry mechanism
// 				return fmt.Errorf("failed to produce kafka message: %w", err)
// 			}
// 		}
// 	}
// 	if order.Status == entity.ORDER_PENDING {
// 		if isAcceptedPayment {
// 			order.SetStatusToOnDelivery()
// 		} else {
// 			order.SetStatusToRejected()
// 		}
// 	}

// 	return nil
// }

func (u *OrderCommandUseCase) UpdateOrderPaymentID(ctx context.Context, order *entity.Order, paymentStatus string) error {
	switch paymentStatus {
	case entity.ORDER_PAYMENT_APPROVED:
		order.SetStatusToOnDelivery()
	case entity.ORDER_PAYMENT_REJECTED:
		order.SetStatusToRejected()
	}

	// if payment rejected, call moveout in warehouse service, put back to warehouse

	return u.repoPostgresCommand.UpdatePaymentID(ctx, order)
}

func (u *OrderCommandUseCase) UpdateOrderStatus(ctx context.Context, order *entity.Order, orderStatus string) error {
	switch orderStatus {
	case entity.ORDER_DELIVERED:
		order.SetStatusToDelivered()
	case entity.ORDER_REJECTED:
		order.SetStatusToRejected()
	}

	return u.repoPostgresCommand.UpdateStatus(ctx, order)
}
