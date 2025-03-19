package app

import "github.com/labstack/echo/v4"

func (h *Handler) Slots(ctx echo.Context) error {
	rows, err := h.db.QueryContext(ctx.Request().Context(), "SELECT slot FROM delivery;")
	if err != nil {
		return internalError(err)
	}

	slots := make([]string, 0, 3)
	for rows.Next() {
		var slot string
		if err := rows.Scan(&slot); err != nil {
			return internalError(err)
		}
		slots = append(slots, slot)
	}

	return ctx.JSON(200, map[string]interface{}{
		"slots": slots,
	})
}
