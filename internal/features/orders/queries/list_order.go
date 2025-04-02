package queries

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/orders"
)

type ListItemOrderDTO struct {
	ID        uuid.UUID `json:"id" path:"id" example:"00000000-0000-0000-0000-000000000000" doc:"Product ID"`
	ProductID uuid.UUID `json:"productId" path:"productId" example:"00000000-0000-0000-0000-000000000000" doc:"Product ID"`
	Quantity  int       `json:"quantity" path:"quantity" example:"1" doc:"Quantity"`
}

type ListOrdersResponse struct {
	Body struct {
		Orders []ListItemOrderDTO `json:"orders" doc:"List of orders"`
	}
}

type ListOrdersHandler func(ctx context.Context) ([]ListItemOrderDTO, error)

func NewListOrdersController(handler ListOrdersHandler) (huma.Operation, func(ctx context.Context, _ *struct{}) (*ListOrdersResponse, error)) {
	return huma.Operation{
			OperationID: "listOrders",
			Method:      http.MethodGet,
			Path:        "/orders",
			Summary:     "List all orders",
			Tags:        []string{"orders"},
		},
		func(ctx context.Context, _ *struct{}) (*ListOrdersResponse, error) {
			products, err := handler(ctx)
			if err != nil {
				return nil, err
			}

			r := &ListOrdersResponse{}
			r.Body.Orders = products
			return r, nil
		}
}

func NewListOrdersHandler(repo orders.Repository) ListOrdersHandler {
	return func(ctx context.Context) ([]ListItemOrderDTO, error) {
		orders, err := repo.ListAll(ctx)
		if err != nil {
			return nil, err
		}

		var dtos []ListItemOrderDTO
		for _, p := range orders {
			dtos = append(dtos, ListItemOrderDTO{
				ID:        p.ID(),
				ProductID: p.ProductID(),
				Quantity:  p.Quantity(),
			})
		}
		return dtos, nil
	}
}
