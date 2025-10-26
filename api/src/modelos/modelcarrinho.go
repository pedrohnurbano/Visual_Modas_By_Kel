package modelos

import "time"

// ItemCarrinho representa um item no carrinho do usuário
type ItemCarrinho struct {
	ID           uint64    `json:"id,omitempty"`
	UsuarioID    uint64    `json:"usuarioId,omitempty"`
	ProdutoID    uint64    `json:"produtoId,omitempty"`
	Quantidade   int       `json:"quantidade"`
	CriadoEm     time.Time `json:"criadoEm,omitempty"`
	AtualizadoEm time.Time `json:"atualizadoEm,omitempty"`
}

// ItemCarrinhoComProduto representa um item do carrinho com os dados completos do produto
type ItemCarrinhoComProduto struct {
	ID           uint64    `json:"id,omitempty"`
	UsuarioID    uint64    `json:"usuarioId,omitempty"`
	ProdutoID    uint64    `json:"produtoId,omitempty"`
	Quantidade   int       `json:"quantidade"`
	CriadoEm     time.Time `json:"criadoEm,omitempty"`
	AtualizadoEm time.Time `json:"atualizadoEm,omitempty"`
	Produto      Produto   `json:"produto"`
}

// ResumoCarrinho representa o resumo do carrinho com totais
type ResumoCarrinho struct {
	Itens           []ItemCarrinhoComProduto `json:"itens"`
	QuantidadeTotal int                      `json:"quantidadeTotal"`
	ValorTotal      float64                  `json:"valorTotal"`
}
