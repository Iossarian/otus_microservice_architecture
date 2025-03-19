package domain

import "github.com/google/uuid"

type Order struct {
	ID     uuid.UUID
	Price  float64
	UserID int
}
