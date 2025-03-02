package app

import (
	"net/http"

	"orchestrator/internal/domain/saga/createorder"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Request struct {
	GoodQty      int     `json:"quantity"`
	Price        float64 `json:"price"`
	DeliverySlot string  `json:"slot"`
}

func (h *Handler) Create(ctx echo.Context) error {
	var req Request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	workflow := createorder.NewWorkflow()
	workflow.AddStep(createorder.SagaStep{
		Name:         "charge",
		Transaction:  h.billing.Charge,
		Compensation: h.billing.Refund,
	})
	workflow.AddStep(createorder.SagaStep{
		Name:         "reserve",
		Transaction:  h.warehouse.Reserve,
		Compensation: h.warehouse.Release,
	})
	workflow.AddStep(createorder.SagaStep{
		Name:         "deliver",
		Transaction:  h.delivery.Reserve,
		Compensation: h.delivery.Cancel,
	})

	if err := workflow.Execute(ctx, createorder.Request{
		UUID:         uuid.New(),
		GoodQty:      req.GoodQty,
		Price:        req.Price,
		DeliverySlot: req.DeliverySlot,
	}); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "order was not created",
		})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"message": "order created",
	})
}
