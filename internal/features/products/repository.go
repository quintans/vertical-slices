package products

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/products/domain"
	"github.com/quintans/vertical-slices/internal/infra"
	"github.com/quintans/vertical-slices/internal/shared/fails"
)

type Repo struct {
	db *infra.DB[*domain.Product]
}

func NewRepository() *Repo {
	return &Repo{
		db: infra.NewDB[*domain.Product](),
	}
}

func (r *Repo) GetByID(_ context.Context, id uuid.UUID) (*domain.Product, error) {
	p, err := r.db.GetByID(id)
	if err != nil {
		if errors.Is(err, infra.ErrDoesNotExist) {
			return nil, fails.ErrNotFound
		}
		return nil, err
	}

	return p, nil
}

func (r *Repo) ListAll(_ context.Context) ([]*domain.Product, error) {
	data := r.db.ListAll()
	return data, nil
}

func (r *Repo) Create(_ context.Context, p *domain.Product) error {
	err := r.db.Create(p.ID(), p)
	if err != nil {
		if errors.Is(err, infra.ErrUniquenessViolation) {
			return fails.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *Repo) Delete(_ context.Context, id uuid.UUID) error {
	r.db.Delete(id)
	return nil
}

func (r *Repo) Update(ctx context.Context, id uuid.UUID, handler func(context.Context, *domain.Product) error) error {
	err := r.db.Update(id, func(p *domain.Product) (*domain.Product, error) {
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

func (r *Repo) GetProductQuantity(_ context.Context, id uuid.UUID) (int, error) {
	p, err := r.db.GetByID(id)
	if err != nil {
		if errors.Is(err, infra.ErrDoesNotExist) {
			return 0, fmt.Errorf("no product with id '%s': %w", id, fails.ErrNotFound)
		}
		return 0, err
	}

	return p.Quantity(), nil
}
