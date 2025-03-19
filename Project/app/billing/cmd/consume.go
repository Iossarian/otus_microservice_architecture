package cmd

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"billing/internal/build"
	"billing/internal/infrastructure/kafka"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func consumeCmd(ctx context.Context, builder *build.Builder) *cobra.Command {
	cmd := &cobra.Command{
		Use: "consume",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("start " + cmd.Short)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("stop " + cmd.Short)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage() //nolint:wrapcheck
		},
	}

	cmd.AddCommand(
		consumeUserRegisteredCmd(ctx, builder),
	)

	return cmd
}

func consumeUserRegisteredCmd(ctx context.Context, builder *build.Builder) *cobra.Command {
	return &cobra.Command{
		Use:   "user-registered",
		Short: "consumer for reading user.registered topic",
		RunE:  runConsumers(ctx, []consumerConstructor{builder.UserCreatedConsumer}, builder),
	}
}

type consumerConstructor func(ctx context.Context) (*kafka.Consumer, error)

func runConsumers(
	ctx context.Context,
	constructors []consumerConstructor,
	builder *build.Builder,
) func(cmd *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		stop := builder.ShutdownChannel(ctx)

		srv, err := builder.HTTPServer(ctx)
		if err != nil {
			return errors.Wrap(err, "build http server")
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Println("run http server")
			}
		}()

		wg := sync.WaitGroup{}

		localCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		for _, construct := range constructors {
			consumer, err := construct(localCtx)
			if err != nil {
				cancel()

				return err
			}

			wg.Add(1)

			go func(consumer *kafka.Consumer) {
				defer wg.Done()

				for {
					select {
					case <-ctx.Done():
						return
					case <-stop:
						return
					default:
						if err := consumer.Consume(ctx); err != nil {
							fmt.Printf("error: %v\n", err)
						}
					}

					time.Sleep(time.Second)
				}
			}(consumer)
		}

		wg.Wait()

		return nil
	}
}
