package app

import (
	"net/http"

	"order/internal/domain"

	"github.com/labstack/echo/v4"
)

func (h *Handler) List(ctx echo.Context) error {
	type Request struct {
		UserID int `json:"user_id"`
	}

	var req Request

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	rows, err := h.db.QueryContext(
		ctx.Request().Context(),
		"SELECT id, price FROM orders WHERE user_id = $1",
		req.UserID,
	)
	if err != nil {
		return internalError(err)
	}

	defer rows.Close()

	orders := make([]domain.Order, 0)

	for rows.Next() {
		var order domain.Order

		if err := rows.Scan(&order.ID, &order.Price); err != nil {
			return internalError(err)
		}
		orders = append(orders, order)
	}

	return ctx.JSON(http.StatusOK, map[string][]domain.Order{
		"orders": orders,
	})
}
