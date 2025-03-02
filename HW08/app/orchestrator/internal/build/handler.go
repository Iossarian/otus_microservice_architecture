package build

import (
	"orchestrator/internal/app"
)

func (b *Builder) handler() (*app.Handler, error) {
	return app.NewHandler(
		b.warehouseClient(),
		b.billingClient(),
		b.deliveryClient(),
	), nil
}
