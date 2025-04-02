package queries

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/products"
)

type GetProductRequest struct {
	ID uuid.UUID `path:"id" doc:"Product ID"`
}

type ProductDTO struct {
	ID    uuid.UUID `json:"id" path:"id" example:"00000000-0000-0000-0000-000000000000" doc:"Product ID"`
	SKU   string    `json:"sku" path:"sku" maxLength:"15" example:"P001" doc:"Product SKU"`
	Name  string    `json:"name" path:"name" maxLength:"30" example:"Product 1" doc:"Product name"`
	Price float64   `json:"price" path:"price" example:"10.99" doc:"Product price"`
}

type GetProductResponse struct {
	Body struct {
		Product ProductDTO `json:"product" doc:"Product"`
	}
}

type GetProductHandler func(ctx context.Context, id uuid.UUID) (*ProductDTO, error)

func NewGetProductController(handler GetProductHandler) (huma.Operation, func(ctx context.Context, input *GetProductRequest) (*GetProductResponse, error)) {
	return huma.Operation{
			OperationID: "getProduct",
			Method:      http.MethodGet,
			Path:        "/products/{id}",
			Summary:     "Get a Product",
			Tags:        []string{"products"},
		},
		func(ctx context.Context, input *GetProductRequest) (*GetProductResponse, error) {
			product, err := handler(ctx, input.ID)
			if err != nil {
				return nil, err
			}

			r := &GetProductResponse{}
			r.Body.Product = *product
			return r, nil
		}
}

func NewGetProductHandler(repo products.Repository) GetProductHandler {
	return func(ctx context.Context, id uuid.UUID) (*ProductDTO, error) {
		product, err := repo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return &ProductDTO{
			ID:    product.ID(),
			SKU:   product.SKU(),
			Name:  product.Name(),
			Price: product.Price(),
		}, nil
	}
}
