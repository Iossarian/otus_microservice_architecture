package app

type Handler struct {
	warehouse Warehouse
	billing   Billing
	delivery  Delivery
}

func NewHandler(
	warehouse Warehouse,
	billing Billing,
	delivery Delivery,
) *Handler {
	return &Handler{
		warehouse: warehouse,
		billing:   billing,
		delivery:  delivery,
	}
}
