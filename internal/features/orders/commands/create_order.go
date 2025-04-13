package commands

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
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

func RegisterCreateOrderController(api huma.API, repo Creater, policy domain.CreateOrderPolicy) {
	handler := NewCreateOrderHandler(repo, policy)

	huma.Register(
		api,
		huma.Operation{
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
		},
	)
}

type Creater interface {
	Create(ctx context.Context, product *domain.Order) error
}

func NewCreateOrderHandler(repo Creater, policy domain.CreateOrderPolicy) func(ctx context.Context, cmd *CreateOrderCommand) (uuid.UUID, error) {
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
