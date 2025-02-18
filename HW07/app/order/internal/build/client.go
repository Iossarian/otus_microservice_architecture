package build

import "order/infrastructure/billing"

func (b *Builder) billingClient() *billing.Client {
	httpClient := b.httpClient()

	return billing.NewClient(
		httpClient,
		b.config.Billing.BaseURL,
	)
}
