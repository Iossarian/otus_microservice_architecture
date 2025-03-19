package app

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Get(ctx echo.Context) error {
	idStr := ctx.Param("id")
	if idStr == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return internalError(err)
	}

	order, err := h.orderRepository.Get(ctx.Request().Context(), id)
	if err != nil {
		return internalError(err)
	}

	price := strconv.Itoa(int(order.Price))

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"price": price,
	})
}
