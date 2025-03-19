package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"delivery/config"

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
						CREATE TABLE IF NOT EXISTS delivery (
						id SERIAL PRIMARY KEY,
						slot VARCHAR,
						created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
						updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
						);
			`)
			if err != nil {
				return errors.Wrap(err, "create delivery table")
			}

			_, err = conn.ExecContext(ctx, `
						INSERT INTO delivery (slot) VALUES ('morning'), ('afternoon'), ('evening');
			`)
			if err != nil {
				return errors.Wrap(err, "insert into delivery table")
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

			fmt.Println("delivery table created")

			return nil
		},
	}
}
