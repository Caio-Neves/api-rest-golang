package errors

import "errors"

//erros da categoria
var (
	ErrCategoriaJaCadastrada         = errors.New("categoria já cadastrada")
	ErrCategoriaNaoCadastrada        = errors.New("categoria não cadastrada")
	ErrNomeCategoriaObrigatorio      = errors.New("nome da categoria deve ser informado")
	ErrDescricaoCategoriaObrigatorio = errors.New("descrição da categoria deve ser informada")
)

//erros específicos
var (
	ErrUuidInvalido         = errors.New("uuid inválido")
	ErrAtributoNaoExistente = errors.New("atributo não existente para atualização")
)
