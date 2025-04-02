package commands

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/orders"
	"github.com/quintans/vertical-slices/internal/features/orders/domain"
)

type CreateOrderCommand struct {
	ProductID uuid.UUID `json:"productId" path:"productId" example:"00000000-0000-0000-0000-000000000000" doc:"Product ID"`
	Quantity  int       `json:"quantity" path:"quantity" example:"1" doc:"Quantity"`
}

type CreateOrderResponse struct {
	Body struct {
		ID uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000" doc:"Order ID"`
	}
}

type CreateOrderHandler func(ctx context.Context, cmd *CreateOrderCommand) (uuid.UUID, error)

func NewCreateOrderController(handler CreateOrderHandler) (huma.Operation, func(ctx context.Context, cmd *CreateOrderCommand) (*CreateOrderResponse, error)) {
	return huma.Operation{
			OperationID: "createOrder",
			Method:      http.MethodPost,
			Path:        "/orders",
			Summary:     "Create Order",
			Description: "Create a new Order for a given product ID and quantity",
			Tags:        []string{"orders"},
		},
		func(ctx context.Context, cmd *CreateOrderCommand) (*CreateOrderResponse, error) {
			id, err := handler(ctx, cmd)
			if err != nil {
				return nil, err
			}

			r := &CreateOrderResponse{}
			r.Body.ID = id
			return r, nil
		}
}

func NewCreateOrderHandler(repo orders.Repository, policy domain.CreateOrderPolicy) CreateOrderHandler {
	return func(ctx context.Context, cmd *CreateOrderCommand) (uuid.UUID, error) {
		p, err := domain.NewOrder(ctx, cmd.ProductID, cmd.Quantity, policy)
		if err != nil {
			return uuid.Nil, err
		}

		err = repo.Create(ctx, p)
		if err != nil {
			return uuid.Nil, err
		}

		return p.ID(), nil
	}
}
