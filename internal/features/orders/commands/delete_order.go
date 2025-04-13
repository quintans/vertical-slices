package commands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type DeleteOrderCommand struct {
	ID uuid.UUID `path:"id" doc:"Order ID"`
}

func RegisterDeleteOrderController(api huma.API, repo Deleter) {
	handler := NewDeleteOrderHandler(repo)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "deleteOrder",
			Method:      http.MethodDelete,
			Path:        "/orders/{id}",
			Summary:     "Delete Order",
			Description: "Delete an order by ID",
			Tags:        []string{"orders"},
		},
		func(ctx context.Context, cmd *DeleteOrderCommand) (*struct{}, error) {
			err := handler(ctx, cmd.ID)

			return nil, err
		},
	)
}

type Deleter interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewDeleteOrderHandler(repo Deleter) func(ctx context.Context, id uuid.UUID) error {
	return func(ctx context.Context, id uuid.UUID) error {
		err := repo.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("deleting order (%s): %w", id, err)
		}

		return nil
	}
}
