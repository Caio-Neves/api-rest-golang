package product

import (
	"context"
	"errors"
	"rest-api-example/category"
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

func (s ProductService) GetAllProducts(ctx context.Context, filters map[string][]string) ([]entities.Product, int, error) {
	op := "ProductService.GetAllProducts()"
	products, totalCount, err := s.productRepository.GetAllProducts(ctx, filters)
	if err != nil {
		return nil, 0, entities.NewInternalServerErrorError(err, op)
	}
	return products, totalCount, nil
}

func (s ProductService) GetProductById(ctx context.Context, id uuid.UUID) (entities.Product, error) {
	op := "ProductService.GetProductById()"
	product, err := s.productRepository.GetProductById(ctx, id)
	if err != nil {
		return entities.Product{}, entities.NewInternalServerErrorError(err, op)
	}
	if product.IsEmpty() {
		return entities.Product{}, entities.NewNotFoundError(ErrProdutoNaoCdastrado, ErrProdutoNaoCdastrado.Error(), op)
	}
	return product, nil
}

func (s ProductService) DeleteProductById(ctx context.Context, id uuid.UUID) error {
	op := "ProductService.DeleteProductById()"
	product, err := s.productRepository.GetProductById(ctx, id)
	if err != nil {
		return entities.NewInternalServerErrorError(err, op)
	}
	if product.IsEmpty() {
		return entities.NewNotFoundError(ErrProdutoNaoCdastrado, ErrProdutoNaoCdastrado.Error(), op)
	}
	err = s.productRepository.DeleteProductById(ctx, id)
	if err != nil {
		return entities.NewInternalServerErrorError(err, op)
	}
	return nil
}

func (s ProductService) DeleteProducts(ctx context.Context, ids []uuid.UUID) error {
	op := "ProductService.DeleteProducts()"
	err := s.productRepository.DeleteProducts(ctx, ids)
	if err != nil {
		entities.NewInternalServerErrorError(err, op)
	}
	return nil
}

func (s ProductService) CreateProduct(ctx context.Context, product entities.Product) (entities.Product, error) {
	op := "ProductService.CreateProcut()"
	if len(product.CategoriesId) == 0 {
		return entities.Product{}, entities.NewBadRequestError(ErrCategoriaDoProdutoEhObrigatoria, ErrCategoriaDoProdutoEhObrigatoria.Error(), op)
	}
	categories, err := s.categoryRepository.GetCategoriesByIds(ctx, product.CategoriesId)
	if err != nil {
		return entities.Product{}, entities.NewInternalServerErrorError(err, op)
	}
	if len(categories) < len(product.CategoriesId) {
		return entities.Product{}, entities.NewBadRequestError(category.ErrCategoriaNaoCadastrada, category.ErrCategoriaNaoCadastrada.Error(), op)

	}
	if product.Name == "" {
		return entities.Product{}, entities.NewBadRequestError(ErrNomeProdutoEhObrigatorio, ErrNomeProdutoEhObrigatorio.Error(), op)

	}
	if product.Description == "" {
		return entities.Product{}, entities.NewBadRequestError(ErrDescricaoProdutoEhObrigatorio, ErrDescricaoProdutoEhObrigatorio.Error(), op)
	}
	product.Id = uuid.New()
	_, err = s.productRepository.CreateProduct(ctx, product)
	if err != nil {
		return entities.Product{}, entities.NewInternalServerErrorError(err, op)
	}
	return product, nil
}

func (s ProductService) UpdateProductFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) (entities.Product, error) {
	op := "ProductService.UpdateProductFields()"

	productDatabase, err := s.productRepository.GetProductById(ctx, id)
	if err != nil {
		return entities.Product{}, entities.NewInternalServerErrorError(err, op)
	}
	if productDatabase.IsEmpty() {
		return entities.Product{}, entities.NewNotFoundError(ErrProdutoNaoCdastrado, ErrProdutoNaoCdastrado.Error(), op)
	}

	product, err := s.productRepository.UpdateProductFields(ctx, id, fields)
	if err != nil {
		return entities.Product{}, entities.NewInternalServerErrorError(err, op)
	}
	return product, nil
}
