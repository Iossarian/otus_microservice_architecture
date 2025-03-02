package build

import (
	"context"

	"billing/config"
)

type Builder struct {
	config   config.Config
	shutdown shutdown
}

func New(_ context.Context, conf config.Config) *Builder {
	b := Builder{config: conf}

	return &b
}
