package modelos

import "time"

// Produto representa um produto no sistema
type Produto struct {
	ID           uint64    `json:"id,omitempty"`
	Nome         string    `json:"nome,omitempty"`
	Descricao    string    `json:"descricao,omitempty"`
	Preco        float64   `json:"preco,omitempty"`
	Tamanho      string    `json:"tamanho,omitempty"`
	Categoria    string    `json:"categoria,omitempty"`
	FotoURL      string    `json:"foto_url,omitempty"`
	UsuarioID    uint64    `json:"usuario_id,omitempty"`
	Ativo        bool      `json:"ativo"`
	CriadoEm     time.Time `json:"criadoEm,omitempty"`
	AtualizadoEm time.Time `json:"atualizadoEm,omitempty"`
}
