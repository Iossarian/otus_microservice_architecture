package build

import (
	"order/internal/infrastructure/billing"
	"order/internal/infrastructure/delivery"
	"order/internal/infrastructure/orchestrator"
	"order/internal/infrastructure/warehouse"
)

func (b *Builder) orchestratorClient() *orchestrator.Client {
	httpClient := b.httpClient()

	return orchestrator.NewClient(
		httpClient,
		b.config.Orchestrator.BaseURL,
	)
}

func (b *Builder) warehouseClient() *warehouse.Client {
	httpClient := b.httpClient()

	return warehouse.NewClient(
		httpClient,
		b.config.Warehouse.BaseURL,
	)
}

func (b *Builder) billingClient() *billing.Client {
	httpClient := b.httpClient()

	return billing.NewClient(
		httpClient,
		b.config.Billing.BaseURL,
	)
}

func (b *Builder) deliveryClient() *delivery.Client {
	httpClient := b.httpClient()

	return delivery.NewClient(
		httpClient,
		b.config.Delivery.BaseURL,
	)
}
