package errors

import "errors"

// erros da categoria
var (
	ErrCategoriaJaCadastrada         = errors.New("categoria já cadastrada")
	ErrCategoriaNaoCadastrada        = errors.New("categoria não cadastrada")
	ErrNomeCategoriaObrigatorio      = errors.New("nome da categoria deve ser informado")
	ErrDescricaoCategoriaObrigatorio = errors.New("descrição da categoria deve ser informada")
)

// erros do produto
var (
	ErrProdutoNaoCdastrado             = errors.New("produto não cadastrada")
	ErrCategoriaDoProdutoEhObrigatoria = errors.New("produto deve ter ao menos 1 categoria")
	ErrNomeProdutoEhObrigatorio        = errors.New("nome do produto deve ser informado")
	ErrDescricaoProdutoEhObrigatorio   = errors.New("descricao do produto deve ser informada")
)

// erros específicos
var (
	ErrUuidInvalido         = errors.New("uuid inválido")
	ErrAtributoNaoExistente = errors.New("atributo não existente para atualização")
)
