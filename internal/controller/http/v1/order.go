package v1

import (
	"net/http"

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
	ID         uuid.UUID            `json:"id"`
	Status     string               `json:"status"`
	TotalPrice float64              `json:"total_price"`
	Items      []itemsOrderResponse `json:"items"`
	Address    addressOrderResponse `json:"address"`
}

type itemsOrderResponse struct {
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Price     float64   `json:"price"`
	Quantity  int64     `json:"quantity"`
	Note      string    `json:"note"`
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
