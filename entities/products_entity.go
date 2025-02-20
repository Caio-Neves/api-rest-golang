package entities

import (
	"context"

	"github.com/google/uuid"
)

type ProductInterface interface {
	GetAllProducts(ctx context.Context, filters map[string][]string) ([]Product, error)
	GetProductById(ctx context.Context, id uuid.UUID) (Product, error)
	DeleteProductById(ctx context.Context, id uuid.UUID)
	DeleteProducts(ctx context.Context, ids []uuid.UUID)
	CreateProduct(ctx context.Context, product Product) (Product, error)
}

type Product struct {
	Id           uuid.UUID   `json:"-"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Price        float64     `json:"price"`
	Active       bool        `json:"active"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at,omitempty"`
	CategoriesId []uuid.UUID `json:"categoriesIds"`
}

func (p *Product) IsEmpty() bool {
	return p.Id == uuid.Nil
}
