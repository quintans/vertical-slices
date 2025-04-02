package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/lib/bus"
	"github.com/quintans/vertical-slices/internal/shared/events"
)

var ErrInsufficientStock = errors.New("insufficient stock")

type Order struct {
	id        uuid.UUID
	productId uuid.UUID
	quantity  int

	events []bus.Message
}

type CreateOrderPolicy interface {
	GetProductQuantity(ctx context.Context, id uuid.UUID) (int, error)
}

func NewOrder(ctx context.Context, productId uuid.UUID, quantity int, policy CreateOrderPolicy) (*Order, error) {
	qty, err := policy.GetProductQuantity(ctx, productId)
	if err != nil {
		return nil, fmt.Errorf("getting stock quantity: %w", err)
	}

	if qty < quantity {
		return nil, fmt.Errorf("creating order: %w", ErrInsufficientStock)
	}

	id := uuid.New()
	return &Order{
		id:        id,
		productId: productId,
		quantity:  quantity,

		events: []bus.Message{
			events.OrderCreated{
				ID:        id,
				ProductID: productId,
				Quantity:  quantity,
			},
		},
	}, nil
}

func (p *Order) ID() uuid.UUID {
	return p.id
}

func (p *Order) ProductID() uuid.UUID {
	return p.productId
}

func (p *Order) Quantity() int {
	return p.quantity
}

func (p *Order) Events() []bus.Message {
	return p.events
}

func HydrateOrder(id, product uuid.UUID, quantity int) *Order {
	return &Order{
		id:        id,
		productId: product,
		quantity:  quantity,
	}
}
