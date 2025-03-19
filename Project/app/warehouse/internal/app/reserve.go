package app

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Reserve(ctx echo.Context) error {
	type Request struct {
		UUID     string `json:"uuid"`
		Quantity int    `json:"quantity"`
	}

	var req Request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	// Поведение для имитации сбоя
	if req.Quantity <= 1 {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "amount must be greater than 1",
		})
	}

	var currentStock int

	err := h.db.QueryRow("SELECT quantity FROM stock WHERE id = $1", 1).Scan(&currentStock)
	if err != nil {
		return internalError(err)
	}

	newStock := currentStock - req.Quantity

	if newStock < 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "insufficient stock",
		})
	}

	_, err = h.db.Exec("UPDATE stock SET quantity = $1, updated_at = $2 WHERE id = $3",
		newStock,
		time.Now(),
		1,
	)
	if err != nil {
		return internalError(err)
	}

	return nil
}
