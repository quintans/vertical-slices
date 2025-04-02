package queries

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/orders"
)

type GetOrderRequest struct {
	ID uuid.UUID `path:"id" doc:"Order ID"`
}

type OrderDTO struct {
	ID        uuid.UUID `json:"id" path:"id" example:"00000000-0000-0000-0000-000000000000" doc:"Product ID"`
	ProductID uuid.UUID `json:"productId" path:"productId" example:"00000000-0000-0000-0000-000000000000" doc:"Product ID"`
	Quantity  int       `json:"quantity" path:"quantity" example:"1" doc:"Quantity"`
}

type GetOrderResponse struct {
	Body struct {
		Order OrderDTO `json:"order" doc:"Order"`
	}
}

type GetOrderHandler func(ctx context.Context, id uuid.UUID) (*OrderDTO, error)

func NewGetOrderController(handler GetOrderHandler) (huma.Operation, func(ctx context.Context, input *GetOrderRequest) (*GetOrderResponse, error)) {
	return huma.Operation{
			OperationID: "getOrder",
			Method:      http.MethodGet,
			Path:        "/orders/{id}",
			Summary:     "Get an Order",
			Tags:        []string{"orders"},
		},
		func(ctx context.Context, input *GetOrderRequest) (*GetOrderResponse, error) {
			order, err := handler(ctx, input.ID)
			if err != nil {
				return nil, err
			}

			r := &GetOrderResponse{}
			r.Body.Order = *order
			return r, nil
		}
}

func NewGetOrderHandler(repo orders.Repository) GetOrderHandler {
	return func(ctx context.Context, id uuid.UUID) (*OrderDTO, error) {
		product, err := repo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return &OrderDTO{
			ID:        product.ID(),
			ProductID: product.ProductID(),
			Quantity:  product.Quantity(),
		}, nil
	}
}
