package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Stock(ctx echo.Context) error {
	var currentStock int

	err := h.db.QueryRow("SELECT quantity FROM stock WHERE id = $1", 1).Scan(&currentStock)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"stock": currentStock,
	})
}
