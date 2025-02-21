package service

import (
	"context"
	"rest-api-example/entities"
	"rest-api-example/errors"
	"rest-api-example/repositories"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepository  repositories.ProductRepositoryPostgres
	categoryRepository repositories.CategoryRepositoryPostgres
}

func NewProductService(p repositories.ProductRepositoryPostgres, c repositories.CategoryRepositoryPostgres) *ProductService {
	return &ProductService{
		productRepository:  p,
		categoryRepository: c,
	}
}

func (s ProductService) GetAllProducts(ctx context.Context, filters map[string][]string) ([]entities.Product, error) {
	products, err := s.productRepository.GetAllProducts(ctx, filters)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s ProductService) GetProductById(ctx context.Context, id uuid.UUID) (entities.Product, error) {
	product, err := s.productRepository.GetProductById(ctx, id)
	if err != nil {
		return entities.Product{}, err
	}
	return product, nil
}

func (s ProductService) DeleteProductById(ctx context.Context, id uuid.UUID) error {
	product, err := s.productRepository.GetProductById(ctx, id)
	if err != nil {
		return err
	}
	if product.IsEmpty() {
		return errors.ErrProdutoNaoCdastrado
	}
	err = s.productRepository.DeleteProductById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s ProductService) DeleteProducts(ctx context.Context, ids []uuid.UUID) error {
	err := s.productRepository.DeleteProducts(ctx, ids)
	return err
}

func (s ProductService) CreateProduct(ctx context.Context, product entities.Product) (entities.Product, error) {
	if len(product.CategoriesId) == 0 {
		return entities.Product{}, errors.ErrCategoriaDoProdutoEhObrigatoria
	}
	categories, err := s.categoryRepository.GetCategoriesByIds(ctx, product.CategoriesId)
	if err != nil {
		return entities.Product{}, err
	}
	if len(categories) < len(product.CategoriesId) {
		return entities.Product{}, errors.ErrCategoriaNaoCadastrada
	}
	if product.Name == "" {
		return entities.Product{}, errors.ErrNomeProdutoEhObrigatorio
	}
	if product.Description == "" {
		return entities.Product{}, errors.ErrDescricaoProdutoEhObrigatorio
	}
	product.Id = uuid.New()
	_, err = s.productRepository.CreateProduct(ctx, product)
	if err != nil {
		return entities.Product{}, err
	}
	return product, nil
}

func (s ProductService) UpdateProductFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) (entities.Product, error) {
	product, err := s.productRepository.UpdateProductFields(ctx, id, fields)
	if err != nil {
		return entities.Product{}, err
	}
	return product, nil
}
