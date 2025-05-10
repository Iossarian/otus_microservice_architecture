package cmd

import (
	"context"

	"notification/config"
	"notification/internal/build"

	"github.com/spf13/cobra"
)

func Run(ctx context.Context, conf config.Config) error {
	root := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	builder := build.New(ctx, conf)

	defer builder.Shutdown(ctx)

	root.AddCommand(
		restCmd(ctx, conf),
		migrateCmd(ctx, conf),
		consumeCmd(ctx, builder),
	)

	return root.ExecuteContext(ctx)
}
