package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Balance(ctx echo.Context) error {
	type Request struct {
		UserId int `json:"user_id"`
	}

	var req Request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	var currentBalance float64
	err := h.db.QueryRow("SELECT balance FROM account WHERE user_id = $1", req.UserId).
		Scan(&currentBalance)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"balance": currentBalance,
	})
}
