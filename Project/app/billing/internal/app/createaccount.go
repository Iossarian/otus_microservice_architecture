package app

import (
	"context"

	"billing/internal/domain"

	"github.com/pkg/errors"
)

func (h *Handler) CreateAccount(ctx context.Context, user domain.User) error {
	_, err := h.db.ExecContext(
		ctx,
		"INSERT INTO account (user_id, balance) VALUES ($1, $2)", user.ID, 100,
	)
	if err != nil {
		return errors.Wrap(err, "create account")
	}

	return nil
}
