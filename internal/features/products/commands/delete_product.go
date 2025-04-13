package commands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type DeleteProductCommand struct {
	ID uuid.UUID `path:"id" doc:"Product ID"`
}

func RegisterDeleteProductController(api huma.API, repo Deleter) {
	handler := NewDeleteProductHandler(repo)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "deleteProduct",
			Method:      http.MethodDelete,
			Path:        "/products/{id}",
			Summary:     "Delete Product",
			Tags:        []string{"products"},
		},
		func(ctx context.Context, cmd *DeleteProductCommand) (*struct{}, error) {
			err := handler(ctx, cmd.ID)

			return nil, err
		},
	)
}

type Deleter interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewDeleteProductHandler(repo Deleter) func(ctx context.Context, id uuid.UUID) error {
	return func(ctx context.Context, id uuid.UUID) error {
		err := repo.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("deleting product (%s): %w", id, err)
		}

		return nil
	}
}
