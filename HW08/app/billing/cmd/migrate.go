package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"billing/config"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func migrateCmd(ctx context.Context, conf config.Config) *cobra.Command {
	command := &cobra.Command{
		Use:       "migrate",
		Short:     "run db migrations",
		ValidArgs: []string{"postgres"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	command.AddCommand(up(ctx, conf))

	return command
}

func up(ctx context.Context, conf config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "up migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := sql.Open("postgres", conf.PostgresDSN())
			if err != nil {
				return errors.Wrap(err, "open connection")
			}

			_, err = conn.ExecContext(ctx, `
						CREATE TABLE IF NOT EXISTS account (
						id SERIAL PRIMARY KEY,
						balance DECIMAL(10, 2),
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
						);
			`)
			if err != nil {
				return errors.Wrap(err, "create billing table")
			}

			_, err = conn.ExecContext(ctx, `
						INSERT INTO account (balance) VALUES (100.00);
			`)
			if err != nil {
				return errors.Wrap(err, "insert into account table")
			}

			_, err = conn.ExecContext(ctx, `
						CREATE TABLE IF NOT EXISTS transaction (
						id UUID,
						status varchar,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
						);
			`)
			if err != nil {
				return errors.Wrap(err, "create transaction table")
			}

			fmt.Println("billing table created")

			return nil
		},
	}
}
