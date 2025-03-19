package app

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (h *Handler) Reserve(ctx echo.Context) error {
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

	var slotID int

	err := h.db.QueryRow("SELECT id FROM delivery WHERE slot = $1", req.Slot).Scan(&slotID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "slot already reserved",
			})
		}

		return internalError(err)
	}

	_, err = h.db.Exec("DELETE FROM delivery WHERE id = $1", slotID)
	if err != nil {
		return internalError(err)
	}

	return nil
}
