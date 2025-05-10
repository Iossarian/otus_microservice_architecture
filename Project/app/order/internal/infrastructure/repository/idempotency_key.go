package repository

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type IdempotencyKeyRepository struct {
	db *memcache.Client
}

func NewIdempotencyKeyRepository(db *memcache.Client) *IdempotencyKeyRepository {
	return &IdempotencyKeyRepository{
		db: db,
	}
}

func (r *IdempotencyKeyRepository) Exists(key uuid.UUID) (bool, error) {
	val, err := r.db.Get(key.String())
	fmt.Println("val", val)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return false, nil
		}

		return false, errors.Wrap(err, "get idempotency key")
	}

	return true, nil
}

func (r *IdempotencyKeyRepository) Set(key uuid.UUID) error {
	err := r.db.Set(&memcache.Item{
		Key:        key.String(),
		Value:      []byte("1"),
		Expiration: 20,
	})
	if err != nil {
		return errors.Wrap(err, "set idempotency key")
	}

	return nil
}
