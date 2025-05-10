package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Balance(ctx echo.Context) error {
	var currentBalance float64

	err := h.db.QueryRow("SELECT balance FROM account WHERE id = $1", 1).Scan(&currentBalance)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"balance": currentBalance,
	})
}
