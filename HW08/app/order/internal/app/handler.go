package app

type Handler struct {
	orchestrator Orchestrator
	warehouse    Warehouse
	billing      Billing
	delivery     Delivery
}

func NewHandler(
	orchestrator Orchestrator,
	billing Billing,
	warehouse Warehouse,
	delivery Delivery,
) *Handler {
	return &Handler{
		orchestrator: orchestrator,
		warehouse:    warehouse,
		billing:      billing,
		delivery:     delivery,
	}
}
