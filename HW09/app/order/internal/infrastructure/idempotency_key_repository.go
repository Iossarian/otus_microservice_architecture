package infrastructure

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type IdempotencyKeyRepository struct {
	db *sql.DB
}

func NewIdempotencyKeyRepository(db *sql.DB) *IdempotencyKeyRepository {
	return &IdempotencyKeyRepository{
		db: db,
	}
}

func (r *IdempotencyKeyRepository) Delete(ctx context.Context, key uuid.UUID) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM idempotency_keys WHERE key = $1",
		key.String(),
	)
	if err != nil {
		return errors.Wrap(err, "delete idempotency key")
	}

	return nil
}
