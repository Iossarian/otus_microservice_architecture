package domain

import "time"

type OrderCreatedEvent struct {
	UserID     string    `json:"user_id"`
	OrderID    string    `json:"order_id"`
	Price      float64   `json:"price"`
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}
