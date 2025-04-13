package queries

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/products/domain"
)

type ListItemProductDTO struct {
	ID    uuid.UUID `json:"id" path:"id" example:"00000000-0000-0000-0000-000000000000" doc:"Product ID"`
	SKU   string    `json:"sku" path:"sku" maxLength:"15" example:"P001" doc:"Product SKU"`
	Name  string    `json:"name" path:"name" maxLength:"30" example:"Product 1" doc:"Product name"`
	Price float64   `json:"price" path:"price" example:"10.99" doc:"Product price"`
}

type ListProductsResponse struct {
	Body struct {
		Products []ListItemProductDTO `json:"products" doc:"List of products"`
	}
}

func RegisterListProductsController(api huma.API, repo Lister) {
	handler := NewListProductsHandler(repo)

	huma.Register(
		api,
		huma.Operation{
			OperationID: "listProducts",
			Method:      http.MethodGet,
			Path:        "/products",
			Summary:     "List all Product",
			Tags:        []string{"products"},
		},
		func(ctx context.Context, _ *struct{}) (*ListProductsResponse, error) {
			products, err := handler(ctx)
			if err != nil {
				return nil, err
			}

			r := &ListProductsResponse{}
			r.Body.Products = products
			return r, nil
		},
	)
}

type Lister interface {
	ListAll(ctx context.Context) ([]*domain.Product, error)
}

func NewListProductsHandler(repo Lister) func(ctx context.Context) ([]ListItemProductDTO, error) {
	return func(ctx context.Context) ([]ListItemProductDTO, error) {
		products, err := repo.ListAll(ctx)
		if err != nil {
			return nil, err
		}

		var dtos []ListItemProductDTO
		for _, p := range products {
			dtos = append(dtos, ListItemProductDTO{
				ID:    p.ID(),
				SKU:   p.SKU(),
				Name:  p.Name(),
				Price: p.Price(),
			})
		}
		return dtos, nil
	}
}
