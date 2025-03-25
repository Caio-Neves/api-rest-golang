package category

import (
	"context"
	"errors"
	"log"
	"rest-api-example/entities"

	"github.com/google/uuid"
)

var (
	ErrCategoriaJaCadastrada         = errors.New("categoria já cadastrada")
	ErrCategoriaNaoCadastrada        = errors.New("categoria não cadastrada")
	ErrNomeCategoriaObrigatorio      = errors.New("nome da categoria deve ser informado")
	ErrDescricaoCategoriaObrigatorio = errors.New("descrição da categoria deve ser informada")
)

type CategoryService struct {
	categoryRepository entities.CategoryInterface
}

func NewCategoryService(r entities.CategoryInterface) CategoryService {
	return CategoryService{
		categoryRepository: r,
	}
}

func (s CategoryService) GetAllCategories(ctx context.Context, page int, limit int, params map[string][]string) ([]entities.Category, error) {
	categories, err := s.categoryRepository.GetPaginateCategories(ctx, page, limit, params)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s CategoryService) GetCategoryById(ctx context.Context, id uuid.UUID) (entities.Category, error) {
	op := "CategoryService.GetCategoryById()"
	category, err := s.categoryRepository.GetCategoryById(ctx, id)
	if err != nil {
		return entities.Category{}, err
	}
	if category.IsEmpty() {
		return entities.Category{}, entities.NewNotFoundError(ErrCategoriaNaoCadastrada, ErrCategoriaNaoCadastrada.Error(), op)
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
		return entities.Category{}, ErrNomeCategoriaObrigatorio
	}
	if category.Description == "" {
		return entities.Category{}, ErrDescricaoCategoriaObrigatorio
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
		return ErrCategoriaNaoCadastrada
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

func (s CategoryService) GetAllProductsByCategory(ctx context.Context, id uuid.UUID) ([]entities.Product, error) {
	products, err := s.categoryRepository.GetAllProductsByCategory(ctx, id)
	if err != nil {
		return nil, err
	}
	return products, nil
}
