package service

import (
	"context"
	"log"
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

func (s CategoryService) GetAllCategories(ctx context.Context, params map[string][]string) ([]entities.Category, error) {
	categories, err := s.categoryRepository.GetAllCategories(ctx, params)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s CategoryService) GetCategoryById(ctx context.Context, id uuid.UUID) (entities.Category, error) {
	category, err := s.categoryRepository.GetCategoryById(ctx, id)
	if err != nil {
		return entities.Category{}, err
	}
	return category, nil
}

func (s CategoryService) GetCategoriesByIds(ctx context.Context, ids []uuid.UUID) ([]entities.Category, error) {
	categories, err := s.categoryRepository.GetCategoriesByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s CategoryService) CreateCategory(ctx context.Context, category entities.Category) (entities.Category, error) {

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

func (s CategoryService) DeleteCategoryById(ctx context.Context, id uuid.UUID) error {
	category, err := s.categoryRepository.GetCategoryById(ctx, id)
	if err != nil {
		return err
	}
	if category.IsEmpty() {
		return errors.ErrCategoriaNaoCadastrada
	}
	err = s.categoryRepository.DeleteCategoryById(ctx, id)
	return err
}

func (s CategoryService) DeleteCategories(ctx context.Context, ids []uuid.UUID) error {
	err := s.categoryRepository.DeleteCategories(ctx, ids)
	return err
}

func (s CategoryService) UpdateCategoryFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) (entities.Category, error) {
	category, err := s.categoryRepository.UpdateCategoryFields(ctx, id, fields)
	if err != nil {
		log.Println(err)
		return entities.Category{}, err
	}
	return category, nil
}
