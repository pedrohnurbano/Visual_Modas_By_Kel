package modelos

import "time"

// Favorito representa um produto favoritado por um usuário
type Favorito struct {
	ID        uint64    `json:"id,omitempty"`
	UsuarioID uint64    `json:"usuarioId,omitempty"`
	ProdutoID uint64    `json:"produtoId,omitempty"`
	CriadoEm  time.Time `json:"criadoEm,omitempty"`
}

// FavoritoComProduto representa um favorito com os dados completos do produto
type FavoritoComProduto struct {
	ID        uint64    `json:"id,omitempty"`
	UsuarioID uint64    `json:"usuarioId,omitempty"`
	ProdutoID uint64    `json:"produtoId,omitempty"`
	CriadoEm  time.Time `json:"criadoEm,omitempty"`
	Produto   Produto   `json:"produto"`
}
