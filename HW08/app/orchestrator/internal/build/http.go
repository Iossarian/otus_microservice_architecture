package build

import (
	"context"
	"net"
	"net/http"
	"time"
)

func (b *Builder) HTTPServer(ctx context.Context) (*http.Server, error) {
	server := http.Server{
		Addr:              b.config.HTTPAddr(),
		ReadHeaderTimeout: time.Millisecond * 25,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	b.shutdown.add(func(ctx context.Context) error {
		return server.Shutdown(ctx)
	})

	return &server, nil
}

func (b *Builder) httpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}
