package app

import (
	"net/http"
	"strconv"

	"gateway/infrastructure/order"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateOrder(ctx echo.Context) error {
	type Request struct {
		Price float64 `form:"price"`
	}

	var p Request

	if err := ctx.Bind(&p); err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to bind withdraw",
			Code:     http.StatusBadRequest,
		}
	}

	tokenString, err := ctx.Request().Cookie("Authorization")
	if err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to get token",
			Code:     http.StatusBadRequest,
		}
	}

	token, err := h.jwtService.Parse(tokenString.Value)
	if err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to parse token",
			Code:     http.StatusBadRequest,
		}
	}

	userID, err := strconv.Atoi(token.ID)
	if err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to convert user ID",
			Code:     http.StatusBadRequest,
		}
	}

	err = h.orderClient.Place(ctx, order.Request{
		UserID: userID,
		Price:  p.Price,
	})
	if err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to place order",
			Code:     http.StatusBadRequest,
		}
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})

}
