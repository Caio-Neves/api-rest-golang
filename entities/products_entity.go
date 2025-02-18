package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ProductInterface interface {
	GetAllProducts(ctx context.Context, filters map[string][]string) ([]Product, error)
	GetProductById(ctx context.Context, id uuid.UUID) (Product, error)
}

type Product struct {
	Id           uuid.UUID   `json:"-"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Price        float64     `json:"price"`
	Active       bool        `json:"active"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at,omitempty"`
	CategoriesId []uuid.UUID `json:"category"`
}

func (p *Product) IsEmpty() bool {
	return p.Id == uuid.Nil
}
