package build

import (
	"notification/internal/app"

	"github.com/pkg/errors"
)

func (b *Builder) handler() (*app.Handler, error) {
	dbConn, err := b.postgres()
	if err != nil {
		return nil, errors.Wrap(err, "build postgres connection")
	}

	return app.NewHandler(
		dbConn,
	), nil
}
