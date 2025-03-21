package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"order/config"

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
						CREATE TABLE IF NOT EXISTS orders (
						id UUID PRIMARY KEY,
						user_id INT,
						price DECIMAL(10, 2),
						created_at TIMESTAMP NOT NULL DEFAULT NOW()                             
						);
			`)
			if err != nil {
				return errors.Wrap(err, "create orders table")
			}

			fmt.Println("orders table created")

			_, err = conn.ExecContext(ctx, `
						CREATE TABLE IF NOT EXISTS idempotency_keys (
						key UUID PRIMARY KEY,
						order_id UUID REFERENCES orders(id),
						created_at TIMESTAMP NOT NULL DEFAULT NOW()                             
						);
			`)
			if err != nil {
				return errors.Wrap(err, "create idempotency keys table")
			}

			fmt.Println("idempotency keys table created")

			return nil
		},
	}
}
