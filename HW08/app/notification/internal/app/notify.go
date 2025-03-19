package app

import (
	"context"

	"notification/internal/domain"
)

func (h *Handler) Notify(ctx context.Context, m domain.Message) error {
	_, err := h.db.ExecContext(
		ctx,
		"INSERT INTO messages (user_id, order_id, price, status) VALUES ($1, $2, $3, $4)",
		m.UserID,
		m.OrderID,
		m.Price,
		m.Status,
	)
	if err != nil {
		return internalError(err)
	}

	return nil
}
