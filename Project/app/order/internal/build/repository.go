package build

import (
	"order/internal/infrastructure/repository"

	"github.com/pkg/errors"
)

func (b *Builder) idempotencyKeyRepository() (*repository.IdempotencyKeyRepository, error) {
	db, err := b.memcache()
	if err != nil {
		return nil, errors.Wrap(err, "build memcache")
	}

	return repository.NewIdempotencyKeyRepository(db), nil
}
