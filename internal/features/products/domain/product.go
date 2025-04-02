package domain

import (
	"errors"

	"github.com/google/uuid"
)

var ErrInsufficientStock = errors.New("insufficient stock")

// Product represents a product in our catalog.
type Product struct {
	id       uuid.UUID
	sku      string
	name     string
	price    float64
	quantity int
}

// NewProduct creates a new product.
func NewProduct(sku, name string, price float64, quantity int) *Product {
	return &Product{
		id:       uuid.New(),
		sku:      sku,
		name:     name,
		price:    price,
		quantity: quantity,
	}
}

// ID returns the product's ID.
func (p *Product) ID() uuid.UUID {
	return p.id
}

// SKU returns the product's SKU.
func (p *Product) SKU() string {
	return p.sku
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) Price() float64 {
	return p.price
}

func (p *Product) Quantity() int {
	return p.quantity
}

func (p *Product) IncreaseStock(quantity int) {
	p.quantity += quantity
}

func (p *Product) DecreaseStock(quantity int) error {
	if p.quantity < quantity {
		return ErrInsufficientStock
	}
	p.quantity -= quantity
	return nil
}

func HydrateProduct(id uuid.UUID, sku, name string, price float64, quantity int) *Product {
	return &Product{
		id:       id,
		sku:      sku,
		name:     name,
		price:    price,
		quantity: quantity,
	}
}
