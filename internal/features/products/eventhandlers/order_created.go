package eventhandlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/products/domain"
	"github.com/quintans/vertical-slices/internal/lib/eventbus"
	"github.com/quintans/vertical-slices/internal/shared/events"
)

type Updater interface {
	Update(ctx context.Context, id uuid.UUID, handler func(context.Context, *domain.Product) error) error
}

func NewOrderCreatedHandler(repo Updater) eventbus.Handler[events.OrderCreated] {
	return func(ctx context.Context, m events.OrderCreated) error {
		return repo.Update(ctx, m.ProductID, func(ctx context.Context, p *domain.Product) error {
			return p.DecreaseStock(m.Quantity)
		})
	}
}
