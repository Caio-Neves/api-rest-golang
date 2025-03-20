package service

import (
	"context"
	"errors"
	"rest-api-example/entities"

	"github.com/google/uuid"
)

// erros do produto
var (
	ErrProdutoNaoCdastrado             = errors.New("produto não cadastrada")
	ErrCategoriaDoProdutoEhObrigatoria = errors.New("produto deve ter ao menos 1 categoria")
	ErrNomeProdutoEhObrigatorio        = errors.New("nome do produto deve ser informado")
	ErrDescricaoProdutoEhObrigatorio   = errors.New("descricao do produto deve ser informada")
)

type ProductService struct {
	productRepository  entities.ProductInterface
	categoryRepository entities.CategoryInterface
}

func NewProductService(p entities.ProductInterface, c entities.CategoryInterface) ProductService {
	return ProductService{
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
		return ErrProdutoNaoCdastrado
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
		return entities.Product{}, ErrCategoriaDoProdutoEhObrigatoria
	}
	categories, err := s.categoryRepository.GetCategoriesByIds(ctx, product.CategoriesId)
	if err != nil {
		return entities.Product{}, err
	}
	if len(categories) < len(product.CategoriesId) {
		return entities.Product{}, ErrCategoriaNaoCadastrada
	}
	if product.Name == "" {
		return entities.Product{}, ErrNomeProdutoEhObrigatorio
	}
	if product.Description == "" {
		return entities.Product{}, ErrDescricaoProdutoEhObrigatorio
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
