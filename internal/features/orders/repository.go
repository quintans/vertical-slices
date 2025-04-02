package orders

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/orders/domain"
	"github.com/quintans/vertical-slices/internal/shared"
	"github.com/quintans/vertical-slices/internal/shared/fails"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	ListAll(ctx context.Context) ([]*domain.Order, error)
	Create(ctx context.Context, product *domain.Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, handler func(context.Context, *domain.Order) error) error
}

type Repo struct {
	data     map[uuid.UUID]domain.Order
	eventBus shared.Publisher
	mutex    sync.RWMutex
}

func NewRepository(eb shared.Publisher) *Repo {
	return &Repo{
		data:     make(map[uuid.UUID]domain.Order),
		eventBus: eb,
	}
}

func (r *Repo) GetByID(_ context.Context, id uuid.UUID) (*domain.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if p, ok := r.data[id]; ok {
		return &p, nil
	}
	return nil, fails.ErrNotFound
}

func (r *Repo) ListAll(_ context.Context) ([]*domain.Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var list []*domain.Order
	for _, v := range r.data {
		list = append(list, &v)
	}
	return list, nil
}

func (r *Repo) Create(ctx context.Context, p *domain.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.data[p.ID()]; ok {
		return errors.New("already exists")
	}
	r.data[p.ID()] = *p

	// publish event
	// This event can be consumed by an in process subscriber or by a message broker.
	// In either case, in a real application, a transaction is needed to guarantee that the event is published only if the save is successful.
	// In the case of a message broker, to guarantee consistency between the save and the publish, we would use the outbox pattern.
	// Note: The context would be the carrier of the transaction.
	r.eventBus.Publish(ctx, p.Events()...)

	return nil
}

func (r *Repo) Delete(_ context.Context, id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.data[id]; !ok {
		return fails.ErrNotFound
	}
	delete(r.data, id)
	return nil
}

func (r *Repo) Update(ctx context.Context, id uuid.UUID, handler func(context.Context, *domain.Order) error) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	p, ok := r.data[id]
	if !ok {
		return fails.ErrNotFound
	}

	if err := handler(ctx, &p); err != nil {
		return err
	}

	return nil
}
