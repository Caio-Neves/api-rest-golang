package entities

import (
	"context"

	"github.com/google/uuid"
)

type CategoryInterface interface {
	GetPaginateCategories(ctx context.Context, page int, limit int, params map[string][]string) ([]Category, int, error)
	GetCategoryById(ctx context.Context, id uuid.UUID) (Category, error)
	GetCategoriesByIds(ctx context.Context, ids []uuid.UUID) ([]Category, error)
	CreateCategory(ctx context.Context, category Category) (Category, error)
	DeleteCategoryById(ctx context.Context, id uuid.UUID) error
	DeleteCategories(ctx context.Context, ids []uuid.UUID) error
	GetAllProductsByCategory(ctx context.Context, id uuid.UUID) ([]Product, error)
	UpdateCategoryFields(ctx context.Context, id uuid.UUID, fields map[string]any) (Category, error)
}

type Category struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at,omitempty"`
}

type CategoryResource struct {
	Category
	Links Hateoas `json:"_meta"`
}

func (c Category) IsEmpty() bool {
	return c.Id == uuid.Nil
}
