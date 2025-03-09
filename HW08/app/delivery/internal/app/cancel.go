package app

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (h *Handler) Cancel(ctx echo.Context) error {
	type Request struct {
		UUID string `json:"uuid"`
		Slot string `json:"slot"`
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

	_, err = h.db.Exec("INSERT INTO delivery VALUES ($1)", req.Slot)
	if err != nil {
		return internalError(err)
	}

	return nil
}
