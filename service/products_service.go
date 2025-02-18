package service

import (
	"context"
	"rest-api-example/entities"
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
