package build

import (
	"order/internal/app"
)

func (b *Builder) handler() (*app.Handler, error) {
	return app.NewHandler(
		b.orchestratorClient(),
		b.billingClient(),
		b.warehouseClient(),
		b.deliveryClient(),
	), nil
}
