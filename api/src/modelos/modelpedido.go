package modelos

import "time"

// Pedido representa um pedido realizado por um usuário
type Pedido struct {
	ID             uint64    `json:"id,omitempty"`
	UsuarioID      uint64    `json:"usuarioId,omitempty"`
	NomeCompleto   string    `json:"nomeCompleto,omitempty"`
	Email          string    `json:"email,omitempty"`
	Telefone       string    `json:"telefone,omitempty"`
	Endereco       string    `json:"endereco,omitempty"`
	Numero         string    `json:"numero,omitempty"`
	Complemento    string    `json:"complemento,omitempty"`
	Bairro         string    `json:"bairro,omitempty"`
	Cidade         string    `json:"cidade,omitempty"`
	Estado         string    `json:"estado,omitempty"`
	CEP            string    `json:"cep,omitempty"`
	FormaPagamento string    `json:"formaPagamento,omitempty"`
	Status         string    `json:"status,omitempty"`
	CodigoRastreio string    `json:"codigoRastreio,omitempty"`
	Total          float64   `json:"total,omitempty"`
	CriadoEm       time.Time `json:"criadoEm,omitempty"`
	AtualizadoEm   time.Time `json:"atualizadoEm,omitempty"`
}

// ItemPedido representa um item dentro de um pedido
type ItemPedido struct {
	ID            uint64    `json:"id,omitempty"`
	PedidoID      uint64    `json:"pedidoId,omitempty"`
	ProdutoID     uint64    `json:"produtoId,omitempty"`
	NomeProduto   string    `json:"nomeProduto,omitempty"`
	PrecoUnitario float64   `json:"precoUnitario,omitempty"`
	Quantidade    int       `json:"quantidade,omitempty"`
	Tamanho       string    `json:"tamanho,omitempty"`
	Subtotal      float64   `json:"subtotal,omitempty"`
	FotoURL       string    `json:"fotoUrl,omitempty"`
	CriadoEm      time.Time `json:"criadoEm,omitempty"`
}

// PedidoCompleto representa um pedido com todos os seus itens
type PedidoCompleto struct {
	Pedido Pedido       `json:"pedido"`
	Itens  []ItemPedido `json:"itens"`
}

// CriarPedidoRequest representa a requisição para criar um pedido
type CriarPedidoRequest struct {
	NomeCompleto   string `json:"nomeCompleto"`
	Email          string `json:"email"`
	Telefone       string `json:"telefone"`
	Endereco       string `json:"endereco"`
	Numero         string `json:"numero"`
	Complemento    string `json:"complemento"`
	Bairro         string `json:"bairro"`
	Cidade         string `json:"cidade"`
	Estado         string `json:"estado"`
	CEP            string `json:"cep"`
	FormaPagamento string `json:"formaPagamento"`
}
