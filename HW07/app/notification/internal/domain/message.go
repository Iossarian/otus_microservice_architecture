package domain

type Message struct {
	UserID  string  `json:"user_id"`
	OrderID string  `json:"order_id"`
	Price   float64 `json:"price"`
	Status  string  `json:"status"`
}
