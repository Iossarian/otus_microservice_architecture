package build

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func (b *Builder) postgres() (*sql.DB, error) {
	db, err := sql.Open("postgres", b.config.PostgresDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	b.shutdown.add(func(_ context.Context) error {
		if err := db.Close(); err != nil {
			return errors.Wrap(err, "close db connection")
		}
		return nil
	})

	return db, nil
}
