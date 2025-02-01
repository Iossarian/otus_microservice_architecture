package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func Run(ctx context.Context, dsn string) error {
	root := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	root.AddCommand(
		migrateCmd(ctx, dsn),
		restCmd(ctx, dsn),
	)

	return root.ExecuteContext(ctx)
}
