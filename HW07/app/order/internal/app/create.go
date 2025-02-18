package app

import (
	"net/http"
	"strconv"
	"time"

	"order/infrastructure/billing"
	"order/internal/domain"

	"github.com/labstack/echo/v4"
)

const (
	orderStatusSuccess = "success"
	orderStatusFailed  = "failed"
)

type Request struct {
	UserId int     `json:"user_id"`
	Price  float64 `json:"price"`
}

func (h *Handler) Create(ctx echo.Context) error {
	var req Request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	_, err := h.billingClient.Withdraw(
		ctx,
		billing.Request{
			UserID: req.UserId,
			Amount: req.Price,
		})
	if err != nil {
		err = h.orderCreatedProducer.Produce(
			ctx.Request().Context(),
			domain.OrderCreatedEvent{
				UserID:     strconv.Itoa(req.UserId),
				Price:      req.Price,
				Status:     orderStatusFailed,
				OccurredAt: time.Now(),
			},
		)
		if err != nil {
			return internalError(err)
		}

		return internalError(err)
	}

	_, err = h.db.Exec("INSERT INTO orders (user_id, price) VALUES ($1, $2)",
		req.UserId,
		req.Price,
	)
	if err != nil {
		return h.rollback(ctx, req)
	}

	var orderID int
	err = h.db.QueryRowContext(
		ctx.Request().Context(),
		"SELECT id FROM orders WHERE user_id = $1 AND price = $2 LIMIT 1",
		req.UserId,
		req.Price,
	).Scan(&orderID)
	if err != nil {
		return h.rollback(ctx, req)
	}

	err = h.orderCreatedProducer.Produce(
		ctx.Request().Context(),
		domain.OrderCreatedEvent{
			UserID:     strconv.Itoa(req.UserId),
			OrderID:    strconv.Itoa(orderID),
			Price:      req.Price,
			Status:     orderStatusSuccess,
			OccurredAt: time.Now(),
		},
	)
	if err != nil {
		return h.rollback(ctx, req)
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "order created",
	})
}

func (h *Handler) rollback(ctx echo.Context, req Request) error {
	err := h.orderCreatedProducer.Produce(
		ctx.Request().Context(),
		domain.OrderCreatedEvent{
			UserID:     strconv.Itoa(req.UserId),
			Price:      req.Price,
			Status:     orderStatusFailed,
			OccurredAt: time.Now(),
		},
	)
	if err != nil {
		return internalError(err)
	}

	_, err = h.billingClient.Deposit(ctx,
		billing.Request{
			UserID: req.UserId,
			Amount: req.Price,
		})
	if err != nil {
		return internalError(err)
	}

	return internalError(err)
}
