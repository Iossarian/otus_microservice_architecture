package app

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Charge(ctx echo.Context) error {
	type Request struct {
		UUID   string  `json:"uuid"`
		Amount float64 `json:"amount"`
	}

	var req Request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	if req.Amount <= 50 {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "amount must be greater than 50",
		})
	}

	var currentBalance float64

	err := h.db.QueryRow("SELECT balance FROM account WHERE id = $1", 1).Scan(&currentBalance)
	if err != nil {
		return internalError(err)
	}

	newBalance := currentBalance - req.Amount

	_, err = h.db.Exec("UPDATE account SET balance = $1, updated_at = $2 WHERE id = $3",
		newBalance,
		time.Now(),
		1,
	)
	if err != nil {
		return internalError(err)
	}

	return nil
}
