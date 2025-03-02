package cmd

import (
	"context"
	"fmt"
	"net/http"

	"orchestrator/config"
	"orchestrator/internal/build"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func restCmd(ctx context.Context, conf config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "rest",
		Short: "start rest server",
		RunE: func(cmd *cobra.Command, args []string) error {
			builder := build.New(ctx, conf)
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			go func() {
				builder.WaitShutdown(ctx)
				cancel()
			}()

			server := builder.RestServer()

			go func() {
				fmt.Println("starting server at ", conf.HTTPAddr())

				if err := server.Start(conf.HTTPAddr()); err != nil && !errors.Is(err, http.ErrServerClosed) {
					server.Logger.Fatal("shutting down the server")
				}
			}()

			<-ctx.Done()

			return nil
		},
	}
}
