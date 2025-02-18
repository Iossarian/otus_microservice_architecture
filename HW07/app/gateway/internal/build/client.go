package build

import (
	"gateway/infrastructure/billing"
	"gateway/infrastructure/notification"
	"gateway/infrastructure/order"
	"gateway/infrastructure/user"
)

func (b *Builder) userClient() *user.Client {
	httpClient := b.httpClient()

	return user.NewClient(
		httpClient,
		b.config.User.BaseURL,
	)
}

func (b *Builder) billingClient() *billing.Client {
	httpClient := b.httpClient()

	return billing.NewClient(
		httpClient,
		b.config.Billing.BaseURL,
	)
}

func (b *Builder) orderClient() *order.Client {
	httpClient := b.httpClient()

	return order.NewClient(
		httpClient,
		b.config.Order.BaseURL,
	)
}

func (b *Builder) notificationClient() *notification.Client {
	httpClient := b.httpClient()

	return notification.NewClient(
		httpClient,
		b.config.Notification.BaseURL,
	)
}
