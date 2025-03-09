package build

import (
	"orchestrator/internal/infrastructure/billing"
	"orchestrator/internal/infrastructure/delivery"
	"orchestrator/internal/infrastructure/warehouse"
)

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
