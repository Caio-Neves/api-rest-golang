package errors

import "errors"

var (
	ErrCategoriaJaCadastrada         = errors.New("categoria já cadastrada")
	ErrCategoriaNaoCadastrada        = errors.New("categoria não cadastrada")
	ErrNomeCategoriaObrigatorio      = errors.New("nome da categoria deve ser informado")
	ErrDescricaoCategoriaObrigatorio = errors.New("descrição da categoria deve ser informada")
)
