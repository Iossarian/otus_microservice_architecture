package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (h *Handler) Refund(ctx echo.Context) error {
	type Request struct {
		UUID   string  `json:"uuid"`
		Amount float64 `json:"amount"`
		UserID int     `json:"user_id"`
	}

	var req Request
	if err := ctx.Bind(&req); err != nil {
		fmt.Println("body", ctx.Request().Body)
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	var transactionStatus string

	err := h.db.QueryRow("SELECT status FROM transaction WHERE id = $1", req.UUID).Scan(&transactionStatus)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return internalError(err)
		}

		_, err = h.db.Exec("INSERT INTO transaction VALUES ($1, $2)", req.UUID, "processing")
		if err != nil {
			return internalError(err)
		}
	}

	if transactionStatus == "processing" {
		return nil
	}

	var currentBalance float64

	err = h.db.QueryRow("SELECT balance FROM account WHERE id = $1", req.UserID).Scan(&currentBalance)
	if err != nil {
		return internalError(err)
	}

	newBalance := currentBalance + req.Amount

	_, err = h.db.Exec("UPDATE account SET balance = $1, updated_at = $2 WHERE id = $3",
		newBalance,
		time.Now(),
		req.UserID,
	)
	if err != nil {
		return internalError(err)
	}

	_, err = h.db.Exec("UPDATE transaction SET status = $1, updated_at = $2 WHERE id = $3",
		"done",
		time.Now(),
		req.UUID,
	)
	if err != nil {
		return internalError(err)
	}

	return nil
}
