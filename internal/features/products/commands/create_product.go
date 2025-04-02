package commands

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/quintans/vertical-slices/internal/features/products"
	"github.com/quintans/vertical-slices/internal/features/products/domain"
)

// CreateProductCommand is a command for creating a product.
type CreateProductCommand struct {
	SKU   string  `json:"sku" path:"sku" maxLength:"15" example:"P001" doc:"Product SKU"`
	Name  string  `json:"name" path:"name" maxLength:"30" example:"Product 1" doc:"Product name"`
	Price float64 `json:"price" path:"price" example:"10.99" doc:"Product price"`
}

type CreateProductResponse struct {
	Body struct {
		ID uuid.UUID `json:"id" example:"00000000-0000-0000-0000-000000000000" doc:"Product ID"`
	}
}

type CreateProductHandler func(ctx context.Context, cmd *CreateProductCommand) (uuid.UUID, error)

func NewCreateProductController(handler CreateProductHandler) (huma.Operation, func(ctx context.Context, cmd *CreateProductCommand) (*CreateProductResponse, error)) {
	return huma.Operation{
			OperationID: "createProduct",
			Method:      http.MethodPost,
			Path:        "/products",
			Summary:     "Create Product",
			Description: "Create a new product",
			Tags:        []string{"products"},
		},
		func(ctx context.Context, cmd *CreateProductCommand) (*CreateProductResponse, error) {
			id, err := handler(ctx, cmd)
			if err != nil {
				return nil, err
			}

			r := &CreateProductResponse{}
			r.Body.ID = id
			return r, nil
		}
}

// NewCreateProductHandler creates a new CreateProductHandler.
func NewCreateProductHandler(repo products.Repository) CreateProductHandler {
	return func(ctx context.Context, cmd *CreateProductCommand) (uuid.UUID, error) {
		p := domain.NewProduct(cmd.SKU, cmd.Name, cmd.Price, 0)

		err := repo.Create(ctx, p)
		if err != nil {
			return uuid.Nil, err
		}

		return p.ID(), nil
	}
}
