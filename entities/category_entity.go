package entities

import (
	"context"

	"github.com/google/uuid"
)

type CategoryInterface interface {
	GetAllCategories(ctx context.Context, params map[string][]string) ([]Category, error)
	GetCategoryById(ctx context.Context, id uuid.UUID) (Category, error)
	CreateCategory(ctx context.Context, category Category) (Category, error)
	DeleteCategoryById(ctx context.Context, id uuid.UUID) error
	DeleteCategories(ctx context.Context, ids []uuid.UUID) error
	// GetAllProductsByCategory
	// UpdateCategories
}

type Category struct {
	Id          uuid.UUID `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at,omitempty"`
}

func (c Category) IsEmpty() bool {
	return c.Id == uuid.Nil
}
