package build

import (
	"order/internal/infrastructure/orchestrator"
)

func (b *Builder) orchestratorClient() *orchestrator.Client {
	httpClient := b.httpClient()

	return orchestrator.NewClient(
		httpClient,
		b.config.Orchestrator.BaseURL,
	)
}
