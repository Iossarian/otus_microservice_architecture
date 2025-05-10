package order

type Request struct {
	GoodQty        int     `json:"quantity"`
	Price          float64 `json:"price"`
	DeliverySlot   string  `json:"slot"`
	UserID         int     `json:"user_id"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type Response struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	ID    int     `json:"id"`
	Price float64 `json:"price"`
}
