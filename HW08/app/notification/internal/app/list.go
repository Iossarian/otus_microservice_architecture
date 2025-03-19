package app

import (
	"net/http"

	"notification/internal/domain"

	"github.com/labstack/echo/v4"
)

func (h *Handler) List(ctx echo.Context) error {
	type Request struct {
		UserID int `json:"user_id"`
	}

	var r Request
	if err := ctx.Bind(&r); err != nil {
		return &echo.HTTPError{
			Internal: err,
			Message:  "failed to bind user",
			Code:     http.StatusBadRequest,
		}
	}

	rows, err := h.db.QueryContext(
		ctx.Request().Context(),
		"SELECT user_id, order_id, price, status FROM messages WHERE user_id = $1",
		r.UserID,
	)
	if err != nil {
		return internalError(err)
	}

	defer rows.Close()

	messages := make([]domain.Message, 0)

	for rows.Next() {
		var message domain.Message

		if err := rows.Scan(&message.UserID, &message.OrderID, &message.Price, &message.Status); err != nil {
			return internalError(err)
		}
		messages = append(messages, message)
	}

	return ctx.JSON(http.StatusOK, map[string][]domain.Message{
		"messages": messages,
	})
}
