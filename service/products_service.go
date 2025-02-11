package service

import "rest-api-example/repositories"

type ProductService struct {
	productRepository *repositories.ProductRepository
}

func NewProductService(r *repositories.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: r,
	}
}
