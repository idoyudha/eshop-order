package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-order/config"
)

const (
	UserIDKey = "userID"
	TokenKey  = "token"
)

type authSuccessResponse struct {
	Code    int          `json:"code"`
	Data    authResponse `json:"data"`
	Message string       `json:"message"`
}

type authResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
}

func cognitoMiddleware(auth config.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, newUnauthorizedError("unauthorized"))
			ctx.Abort()
			return
		}

		tokenString = strings.TrimSpace(strings.Replace(tokenString, "Bearer ", "", 1))
		ctx.Set(TokenKey, tokenString)

		authURL := fmt.Sprintf("%s/v1/auth/%s", auth.BaseURL, tokenString)
		response, err := http.Get(authURL)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
			ctx.Abort()
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			ctx.JSON(http.StatusUnauthorized, newUnauthorizedError("unauthorized"))
			ctx.Abort()
			return
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
			ctx.Abort()
			return
		}

		var authSuccessResponse authSuccessResponse
		if err := json.Unmarshal(body, &authSuccessResponse); err != nil {
			ctx.JSON(http.StatusInternalServerError, newInternalServerError(err.Error()))
			ctx.Abort()
			return
		}

		ctx.Set(UserIDKey, authSuccessResponse.Data.UserID)
		ctx.Next()
	}
}
