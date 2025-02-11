package service

import (
	"context"
	"rest-api-example/entities"
	"rest-api-example/errors"

	"github.com/google/uuid"
)

type CategoryService struct {
	categoryRepository entities.CategoryInterface
}

func NewCategoryService(r entities.CategoryInterface) *CategoryService {
	return &CategoryService{
		categoryRepository: r,
	}
}

func (s *CategoryService) GetAllCategories(ctx context.Context, params map[string][]string) ([]entities.Category, error) {
	categories, err := s.categoryRepository.GetAllCategories(ctx, params)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryService) GetCategoryById(ctx context.Context, id string) (entities.Category, error) {
	category, err := s.categoryRepository.GetCategoryById(ctx, id)
	if err != nil {
		return entities.Category{}, err
	}
	return category, nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, category entities.Category) (entities.Category, error) {

	if category.Name == "" {
		return entities.Category{}, errors.ErrNomeCategoriaObrigatorio
	}
	if category.Description == "" {
		return entities.Category{}, errors.ErrDescricaoCategoriaObrigatorio
	}
	category.Id = uuid.New()
	category, err := s.categoryRepository.CreateCategory(ctx, category)
	if err != nil {
		return entities.Category{}, err
	}
	return category, nil
}
