package cmd

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

func migrateCmd(ctx context.Context, dsn string) *cobra.Command {
	command := &cobra.Command{
		Use:       "migrate",
		Short:     "run db migrations",
		ValidArgs: []string{"postgres"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	command.AddCommand(up(ctx, dsn))

	return command
}

func up(ctx context.Context, dsn string) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "up migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := sql.Open("postgres", dsn)
			if err != nil {
				return errors.Wrap(err, "open connection")
			}

			_, err = conn.ExecContext(ctx, `
						CREATE TABLE IF NOT EXISTS users (
						id SERIAL PRIMARY KEY,
						name VARCHAR(100),
						age INT
						);
			`)
			if err != nil {
				return errors.Wrap(err, "create users table")
			}

			return nil
		},
	}
}
