package modelos

import (
	"errors"
	"strings"
	"time"
)

// Produto representa um produto no sistema
type Produto struct {
	ID          uint64    `json:"id,omitempty"`
	Nome        string    `json:"nome,omitempty"`
	Descricao   string    `json:"descricao,omitempty"`
	Preco       float64   `json:"preco,omitempty"`
	Tamanho     string    `json:"tamanho,omitempty"`
	Categoria   string    `json:"categoria,omitempty"`
	FotoURL     string    `json:"foto_url,omitempty"`
	UsuarioID   uint64    `json:"usuario_id,omitempty"`
	Ativo       bool      `json:"ativo"`
	CriadoEm    time.Time `json:"criadoEm,omitempty"`
	AtualizadoEm time.Time `json:"atualizadoEm,omitempty"`
}

// TamanhosValidos lista todos os tamanhos permitidos
var TamanhosValidos = []string{"PP", "P", "M", "G", "GG", "XG", "XGG"}

// CategoriasValidas lista todas as categorias permitidas
var CategoriasValidas = []string{
	"Vestidos",
	"Blusas e Camisas",
	"Calças",
	"Saias",
	"Shorts e Bermudas",
	"Jaquetas e Casacos",
	"Macacões",
	"Blazers",
	"Body",
	"Regatas",
}

// Preparar valida e formata o produto
func (produto *Produto) Preparar(etapa string) error {
	if erro := produto.validar(etapa); erro != nil {
		return erro
	}

	produto.formatar()
	return nil
}

func (produto *Produto) validar(etapa string) error {
	if produto.Nome == "" {
		return errors.New("o nome é obrigatório e não pode estar em branco")
	}

	if len(produto.Nome) > 255 {
		return errors.New("o nome não pode ter mais de 255 caracteres")
	}

	if produto.Descricao == "" {
		return errors.New("a descrição é obrigatória e não pode estar em branco")
	}

	if produto.Preco <= 0 {
		return errors.New("o preço deve ser maior que zero")
	}

	if produto.Tamanho == "" {
		return errors.New("o tamanho é obrigatório")
	}

	// Validar tamanho
	tamanhoValido := false
	for _, t := range TamanhosValidos {
		if produto.Tamanho == t {
			tamanhoValido = true
			break
		}
	}
	if !tamanhoValido {
		return errors.New("tamanho inválido. Use: PP, P, M, G, GG, XG ou XGG")
	}

	if produto.Categoria == "" {
		return errors.New("a categoria é obrigatória")
	}

	// Validar categoria
	categoriaValida := false
	for _, c := range CategoriasValidas {
		if produto.Categoria == c {
			categoriaValida = true
			break
		}
	}
	if !categoriaValida {
		return errors.New("categoria inválida")
	}

	if etapa == "cadastro" && produto.UsuarioID == 0 {
		return errors.New("o ID do usuário é obrigatório")
	}

	return nil
}

func (produto *Produto) formatar() {
	produto.Nome = strings.TrimSpace(produto.Nome)
	produto.Descricao = strings.TrimSpace(produto.Descricao)
	produto.Tamanho = strings.TrimSpace(produto.Tamanho)
	produto.Categoria = strings.TrimSpace(produto.Categoria)
	produto.FotoURL = strings.TrimSpace(produto.FotoURL)
}
