package cmd

import (
	"context"

	"github.com/Iossarian/otus_microservice_architecture/helm/rest"
	"github.com/spf13/cobra"
)

func restCmd(ctx context.Context, dsn string) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "rest",
		Short: "start rest server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			go func() {
				cancel()
			}()

			rest.ListenAndServer(ctx, dsn)

			<-ctx.Done()

			return nil
		},
	}
}
