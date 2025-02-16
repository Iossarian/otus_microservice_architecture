package app

import (
	"net/http"

	"gateway/infrastructure/billing"
	"gateway/infrastructure/notification"
	"gateway/infrastructure/order"
	"gateway/infrastructure/user"
	"gateway/internal/jwt"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	userClient         *user.Client
	billingClient      *billing.Client
	orderClient        *order.Client
	notificationClient *notification.Client
	jwtService         *jwt.Service
}

func NewHandler(
	userClient *user.Client,
	billingClient *billing.Client,
	orderClient *order.Client,
	notificationClient *notification.Client,
	jwtService *jwt.Service,
) *Handler {
	return &Handler{
		userClient:         userClient,
		billingClient:      billingClient,
		orderClient:        orderClient,
		notificationClient: notificationClient,
		jwtService:         jwtService,
	}
}

func internalError(err error) error {
	return &echo.HTTPError{
		Internal: err,
		Message:  "Internal error",
		Code:     http.StatusInternalServerError,
	}
}
