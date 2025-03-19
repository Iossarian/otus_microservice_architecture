package app

import (
	"net/http"
	"strconv"

	"gateway/infrastructure/order"

	"github.com/labstack/echo/v4"
)

type Request struct {
	GoodQty      int     `form:"quantity"`
	Price        float64 `form:"price"`
	DeliverySlot string  `form:"slot"`
}

func (h *Handler) CreateOrder(ctx echo.Context) error {
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
		UserID:       userID,
		Price:        p.Price,
		GoodQty:      p.GoodQty,
		DeliverySlot: p.DeliverySlot,
	})
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})

}
