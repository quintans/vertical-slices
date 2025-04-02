package eventhandlers

import (
	"context"

	"github.com/quintans/vertical-slices/internal/features/products"
	"github.com/quintans/vertical-slices/internal/features/products/domain"
	"github.com/quintans/vertical-slices/internal/lib/bus"
	"github.com/quintans/vertical-slices/internal/shared/events"
)

func NewOrderCreatedHandler(repo products.Repository) bus.Handler[events.OrderCreated] {
	return func(ctx context.Context, m events.OrderCreated) error {
		return repo.Update(ctx, m.ProductID, func(ctx context.Context, p *domain.Product) error {
			return p.DecreaseStock(m.Quantity)
		})
	}
}
