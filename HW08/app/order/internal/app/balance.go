package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Balance(ctx echo.Context) error {
	balance, err := h.billing.Balance(ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "balance error",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"balance": balance,
	})
}
