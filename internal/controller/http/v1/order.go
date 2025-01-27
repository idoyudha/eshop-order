package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/internal/usecase"
	"github.com/idoyudha/eshop-order/pkg/logger"
)

type orderRoutes struct {
	uoc usecase.OrderCommand
	uoq usecase.OrderQuery
	l   logger.Interface
}

func newOrderRoutes(
	handler *gin.RouterGroup,
	uoc usecase.OrderCommand,
	uoq usecase.OrderQuery,
	l logger.Interface,
	authMid gin.HandlerFunc,
) {
	r := &orderRoutes{uoc: uoc, uoq: uoq, l: l}

	h := handler.Group("/orders").Use(authMid)
	{
		h.POST("", r.createOrder)
		h.GET("/user", r.getOrderByUserID)
		h.GET("/:id", r.getOrderByID)
		h.GET("", r.getAllOrders)
		h.PATCH("/:id/status", r.updateOrderStatus)
		h.GET("/:id/ttl", r.getOrderTTL)
	}
}

type createOrderRequest struct {
	Items   []createItemsOrderRequest `json:"items"`
	Address createAddressOrderRequest `json:"address"`
}

type createItemsOrderRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	Price     float64   `json:"price"`
}

type createAddressOrderRequest struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipcode"`
	Note    string `json:"note"`
}

type orderResponse struct {
	ID              uuid.UUID            `json:"id"`
	Status          string               `json:"status"`
	TotalPrice      float64              `json:"total_price"`
	PaymentID       uuid.UUID            `json:"payment_id"`
	PaymentStatus   string               `json:"payment_status"`
	PaymentImageURL string               `json:"payment_image_url"`
	Items           []itemsOrderResponse `json:"items"`
	Address         addressOrderResponse `json:"address"`
	CreatedAt       time.Time            `json:"created_at"`
}

type itemsOrderResponse struct {
	OrderID      uuid.UUID `json:"order_id"`
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ImageURL     string    `json:"image_url"`
	Price        float64   `json:"price"`
	Quantity     int64     `json:"quantity"`
	ShippingCost float64   `json:"shipping_cost"`
	Note         string    `json:"note"`
}

type addressOrderResponse struct {
	OrderID uuid.UUID `json:"order_id"`
	Street  string    `json:"street"`
	City    string    `json:"city"`
	State   string    `json:"state"`
	ZipCode string    `json:"zipcode"`
	Note    string    `json:"note"`
}

func (r *orderRoutes) createOrder(ctx *gin.Context) {
	var req createOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - orderRoutes - createOrder")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	userID, exist := ctx.Get(UserIDKey)
	if !exist {
		r.l.Error("not exist", "http - v1 - orderRoutes - createOrder")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError("user id not exist"))
		return
	}

	token, exist := ctx.Get(TokenKey)
	if !exist {
		r.l.Error("not exist", "http - v1 - orderRoutes - createOrder")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError("token not exist"))
		return
	}

	order, err := CreateOrderRequestToOrderEntity(req, userID.(uuid.UUID))
	if err != nil {
		r.l.Error(err, "http - v1 - orderRoutes - createOrder")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	err = r.uoc.CreateOrder(ctx.Request.Context(), &order, token.(string))
	if err != nil {
		r.l.Error(err, "http - v1 - orderRoutes - createOrder")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	response := OrderEntityToCreatedOrderResponse(order)

	ctx.JSON(http.StatusCreated, newCreateSuccess(response))
}

func (r *orderRoutes) getOrderByID(ctx *gin.Context) {
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - cartRoutes - updateCart")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	order, err := r.uoq.GetOrderByID(ctx.Request.Context(), orderID)
	if err != nil {
		r.l.Error(err, "http - v1 - orderRoutes - getOrderByID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	response := OrderViewEntityToGetOneOrderResponse(order)

	ctx.JSON(http.StatusOK, newGetSuccess(response))
}

func (r *orderRoutes) getOrderByUserID(ctx *gin.Context) {
	userID, exist := ctx.Get(UserIDKey)
	if !exist {
		r.l.Error("not exist", "http - v1 - orderRoutes - getOrderByUserID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError("user id not exist"))
		return
	}

	orders, err := r.uoq.GetOrderByUserID(ctx.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		r.l.Error(err, "http - v1 - orderRoutes - getOrderByUserID")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	response := OrderViewEntityToGetManyOrderResponse(orders)

	ctx.JSON(http.StatusOK, newGetSuccess(response))
}

func (r *orderRoutes) getAllOrders(ctx *gin.Context) {
	orders, err := r.uoq.GetAllOrders(ctx.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - orderRoutes - getAllOrders")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	response := OrderViewEntityToGetManyOrderResponse(orders)

	ctx.JSON(http.StatusOK, newGetSuccess(response))
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

func (r *orderRoutes) updateOrderStatus(ctx *gin.Context) {
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - orderRoutes - updateOrderStatus")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	var req UpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - orderRoutes - updateOrderStatus")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	orderEntity := UpdateOrderRequestToOrderEntity(orderID)

	err = r.uoc.UpdateOrderStatus(ctx.Request.Context(), &orderEntity, req.Status)
	if err != nil {
		r.l.Error(err, "http - v1 - orderRoutes - updateOrderStatus")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, newUpdateSuccess(nil))
}

type orderTTLResponse struct {
	TTL time.Duration `json:"ttl_seconds"`
}

func (r *orderRoutes) getOrderTTL(ctx *gin.Context) {
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		r.l.Error(err, "http - v1 - orderRoutes - getOrderTTL")
		ctx.JSON(http.StatusBadRequest, newBadRequestError(err.Error()))
		return
	}

	ttl, err := r.uoc.GetOrderTTL(ctx.Request.Context(), orderID)
	if err != nil {
		r.l.Error(err, "http - v1 - orderRoutes - getOrderTTL")
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
		return
	}

	var response orderTTLResponse
	response.TTL = ttl

	ctx.JSON(http.StatusOK, newGetSuccess(response))
}
