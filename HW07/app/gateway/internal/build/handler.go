package build

import (
	"gateway/internal/app"
)

func (b *Builder) handler() (*app.Handler, error) {
	return app.NewHandler(
		b.userClient(),
		b.billingClient(),
		b.orderClient(),
		b.notificationClient(),
		b.jwtService(),
	), nil
}
