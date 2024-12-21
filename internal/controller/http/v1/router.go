package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-order/config"
	"github.com/idoyudha/eshop-order/internal/usecase"
	"github.com/idoyudha/eshop-order/pkg/logger"
)

func NewRouter(
	handler *gin.Engine,
	ucq usecase.OrderQuery,
	uoc usecase.OrderCommand,
	l logger.Interface,
	auth config.AuthService,
) {
	handler.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	authMid := cognitoMiddleware(auth)

	h := handler.Group("/v1")
	{
		newOrderRoutes(h, uoc, ucq, l, authMid)
	}
}
