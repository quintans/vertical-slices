package products

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/products/domain"
	"github.com/quintans/vertical-slices/internal/shared/fails"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListAll(ctx context.Context) ([]*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, handler func(context.Context, *domain.Product) error) error
}

type Repo struct {
	data  map[uuid.UUID]domain.Product
	mutex sync.RWMutex
}

func NewRepository() *Repo {
	return &Repo{
		data: make(map[uuid.UUID]domain.Product),
	}
}

func (r *Repo) GetByID(_ context.Context, id uuid.UUID) (*domain.Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if p, ok := r.data[id]; ok {
		return &p, nil
	}
	return nil, fails.ErrNotFound
}

func (r *Repo) ListAll(_ context.Context) ([]*domain.Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var list []*domain.Product
	for _, v := range r.data {
		list = append(list, &v)
	}
	return list, nil
}

func (r *Repo) Create(_ context.Context, p *domain.Product) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.data[p.ID()]; ok {
		return errors.New("already exists")
	}
	r.data[p.ID()] = *p
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

func (r *Repo) Update(ctx context.Context, id uuid.UUID, handler func(context.Context, *domain.Product) error) error {
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

func (r *Repo) Seed(products ...*domain.Product) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, p := range products {
		r.data[p.ID()] = *p
	}
}

func (r *Repo) GetProductQuantity(_ context.Context, id uuid.UUID) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if p, ok := r.data[id]; ok {
		return p.Quantity(), nil
	}
	return 0, fmt.Errorf("no product with id '%s': %w", id, fails.ErrNotFound)
}
