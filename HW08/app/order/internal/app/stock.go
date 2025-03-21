package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Stock(ctx echo.Context) error {
	stock, err := h.warehouse.Stock(ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "stock error",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"stock": stock,
	})
}
