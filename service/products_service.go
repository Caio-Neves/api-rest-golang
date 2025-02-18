package service

import (
	"context"
	"rest-api-example/entities"
	"rest-api-example/errors"
	"rest-api-example/repositories"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepository *repositories.ProductRepositoryPostgres
}

func NewProductService(r *repositories.ProductRepositoryPostgres) *ProductService {
	return &ProductService{
		productRepository: r,
	}
}

func (s *ProductService) GetAllProducts(ctx context.Context, filters map[string][]string) ([]entities.Product, error) {
	products, err := s.productRepository.GetAllProducts(ctx, filters)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProductById(ctx context.Context, id uuid.UUID) (entities.Product, error) {
	product, err := s.productRepository.GetProductById(ctx, id)
	if err != nil {
		return entities.Product{}, err
	}
	return product, nil
}

func (s *ProductService) DeleteProductById(ctx context.Context, id uuid.UUID) error {
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

func (s *ProductService) DeleteProducts(ctx context.Context, ids []uuid.UUID) error {
	err := s.productRepository.DeleteProducts(ctx, ids)
	return err
}
