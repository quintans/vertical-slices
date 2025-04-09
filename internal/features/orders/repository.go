package orders

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/orders/domain"
	"github.com/quintans/vertical-slices/internal/infra"
	"github.com/quintans/vertical-slices/internal/shared"
	"github.com/quintans/vertical-slices/internal/shared/fails"
)

type Repo struct {
	db       *infra.DB[*domain.Order]
	eventBus shared.Publisher
}

func NewRepository(eb shared.Publisher) *Repo {
	return &Repo{
		db:       infra.NewDB[*domain.Order](),
		eventBus: eb,
	}
}

func (r *Repo) GetByID(_ context.Context, id uuid.UUID) (*domain.Order, error) {
	o, err := r.db.GetByID(id)
	if err != nil {
		if errors.Is(err, infra.ErrDoesNotExist) {
			return nil, fails.ErrNotFound
		}
		return nil, err
	}

	return o, nil
}

func (r *Repo) ListAll(_ context.Context) ([]*domain.Order, error) {
	data := r.db.ListAll()
	return data, nil
}

func (r *Repo) Create(ctx context.Context, o *domain.Order) error {
	err := r.db.Create(o.ID(), o)
	if err != nil {
		if errors.Is(err, infra.ErrUniquenessViolation) {
			return fails.ErrAlreadyExists
		}
		return err
	}

	// publish event
	// This event can be consumed by an in process subscriber or by a message broker.
	// In either case, in a real application, a transaction is needed to guarantee that the event is published only if the save is successful.
	// In the case of a message broker, to guarantee consistency between the save and the publish, we would use the outbox pattern.
	// Note: The context would be the carrier of the transaction.
	r.eventBus.Publish(ctx, o.Events()...)

	return nil
}

func (r *Repo) Delete(_ context.Context, id uuid.UUID) error {
	r.db.Delete(id)
	return nil
}

func (r *Repo) Update(ctx context.Context, id uuid.UUID, handler func(context.Context, *domain.Order) error) error {
	err := r.db.Update(id, func(p *domain.Order) (*domain.Order, error) {
		err := handler(ctx, p)
		return p, err
	})
	if err != nil {
		if errors.Is(err, infra.ErrDoesNotExist) {
			return fails.ErrNotFound
		}
		return err
	}

	return nil
}
