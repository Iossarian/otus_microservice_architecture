package createorder

import "github.com/google/uuid"

type Request struct {
	UUID         uuid.UUID
	GoodQty      int
	Price        float64
	DeliverySlot string
	UserID       int
}
