package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Slots(ctx echo.Context) error {
	slots, err := h.delivery.Slots(ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "slots error",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"slots": slots,
	})
}
