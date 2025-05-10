package build

import (
	"order/internal/app"

	"github.com/pkg/errors"
)

func (b *Builder) handler() (*app.Handler, error) {
	dbConn, err := b.postgres()
	if err != nil {
		return nil, errors.Wrap(err, "build postgres connection")
	}

	idempotencyKeyRepository, err := b.idempotencyKeyRepository()
	if err != nil {
		return nil, errors.Wrap(err, "build idempotency key repository")
	}

	return app.NewHandler(
		dbConn,
		b.orchestratorClient(),
		b.OrderCreatedProducer(),
		idempotencyKeyRepository,
	), nil
}
