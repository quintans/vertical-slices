package queries

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/orders/domain"
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

func RegisterGetOrderController(api huma.API, repo Getter) {
	handler := NewGetOrderHandler(repo)

	huma.Register(
		api,
		huma.Operation{
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
		},
	)
}

type Getter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
}

func NewGetOrderHandler(repo Getter) func(ctx context.Context, id uuid.UUID) (*OrderDTO, error) {
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
