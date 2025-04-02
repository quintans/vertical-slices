package events

import "github.com/google/uuid"

type OrderCreated struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Quantity  int
}

func (e OrderCreated) Kind() string {
	return "OrderCreated"
}
