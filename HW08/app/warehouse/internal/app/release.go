package app

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (h *Handler) Release(ctx echo.Context) error {
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

	var currentStock int

	err = h.db.QueryRow("SELECT quantity FROM stock WHERE id = $1", 1).Scan(&currentStock)
	if err != nil {
		return internalError(err)
	}

	newStock := currentStock + req.Quantity

	_, err = h.db.Exec("UPDATE stock SET quantity = $1, updated_at = $2 WHERE id = $3",
		newStock,
		time.Now(),
		1,
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
