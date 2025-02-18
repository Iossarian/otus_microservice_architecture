package app

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Deposit(ctx echo.Context) error {
	type Request struct {
		UserId int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	var req Request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	var currentBalance float64

	err := h.db.QueryRow("SELECT balance FROM account WHERE user_id = $1", req.UserId).Scan(&currentBalance)
	if err != nil {
		return internalError(err)
	}

	newBalance := currentBalance + req.Amount

	_, err = h.db.Exec("UPDATE account SET balance = $1, updated_at = $2 WHERE user_id = $3",
		newBalance,
		time.Now(),
		req.UserId,
	)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"balance": newBalance,
	})
}
