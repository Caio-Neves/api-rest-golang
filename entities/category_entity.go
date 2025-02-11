package entities

import (
	"context"

	"github.com/google/uuid"
)

type CategoryInterface interface {
	GetAllCategories(ctx context.Context, params map[string][]string) ([]Category, error)
	GetCategoryById(ctx context.Context, id string) (Category, error)
	// GetAllProductsByCategory
	CreateCategory(ctx context.Context, category Category) (Category, error)
	// DeleteCategories
	// UpdateCategories
	// DeleteCategoryById
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
