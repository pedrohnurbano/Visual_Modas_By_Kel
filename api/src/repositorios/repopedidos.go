package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

// Pedidos representa um repositório de pedidos
type Pedidos struct {
	db *sql.DB
}

// NovoRepositorioDePedidos cria um repositório de pedidos
func NovoRepositorioDePedidos(db *sql.DB) *Pedidos {
	return &Pedidos{db}
}

// Criar cria um novo pedido e seus itens a partir do carrinho do usuário
func (repositorio Pedidos) Criar(pedido modelos.Pedido) (uint64, error) {
	// Iniciar transação
	tx, erro := repositorio.db.Begin()
	if erro != nil {
		return 0, erro
	}
	defer tx.Rollback()

	// Inserir pedido
	resultado, erro := tx.Exec(`
		INSERT INTO pedidos (
			usuario_id, nome_completo, email, telefone,
			endereco, numero, complemento, bairro, cidade, estado, cep,
			forma_pagamento, status, total
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		pedido.UsuarioID,
		pedido.NomeCompleto,
		pedido.Email,
		pedido.Telefone,
		pedido.Endereco,
		pedido.Numero,
		pedido.Complemento,
		pedido.Bairro,
		pedido.Cidade,
		pedido.Estado,
		pedido.CEP,
		pedido.FormaPagamento,
		pedido.Status,
		pedido.Total,
	)

	if erro != nil {
		return 0, erro
	}

	pedidoID, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	// Buscar itens do carrinho do usuário
	linhas, erro := tx.Query(`
		SELECT 
			c.produto_id, c.quantidade,
			p.nome, p.preco, p.tamanho, p.foto_url
		FROM carrinho c
		INNER JOIN produtos p ON c.produto_id = p.id
		WHERE c.usuario_id = ? AND p.ativo = true
	`, pedido.UsuarioID)

	if erro != nil {
		return 0, erro
	}
	defer linhas.Close()

	// Inserir itens do pedido
	for linhas.Next() {
		var produtoID uint64
		var quantidade int
		var nome, tamanho, fotoURL string
		var preco float64

		if erro = linhas.Scan(&produtoID, &quantidade, &nome, &preco, &tamanho, &fotoURL); erro != nil {
			return 0, erro
		}

		subtotal := preco * float64(quantidade)

		_, erro = tx.Exec(`
			INSERT INTO itens_pedido (
				pedido_id, produto_id, nome_produto, preco_unitario,
				quantidade, tamanho, subtotal, foto_url
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, pedidoID, produtoID, nome, preco, quantidade, tamanho, subtotal, fotoURL)

		if erro != nil {
			return 0, erro
		}
	}

	// Limpar carrinho do usuário
	_, erro = tx.Exec("DELETE FROM carrinho WHERE usuario_id = ?", pedido.UsuarioID)
	if erro != nil {
		return 0, erro
	}

	// Commit da transação
	if erro = tx.Commit(); erro != nil {
		return 0, erro
	}

	return uint64(pedidoID), nil
}

// BuscarPorID busca um pedido pelo ID
func (repositorio Pedidos) BuscarPorID(pedidoID uint64) (modelos.Pedido, error) {
	linha := repositorio.db.QueryRow(`
		SELECT 
			id, usuario_id, nome_completo, email, telefone,
			endereco, numero, complemento, bairro, cidade, estado, cep,
			forma_pagamento, status, total, criadoEm, atualizadoEm
		FROM pedidos
		WHERE id = ?
	`, pedidoID)

	var pedido modelos.Pedido
	erro := linha.Scan(
		&pedido.ID,
		&pedido.UsuarioID,
		&pedido.NomeCompleto,
		&pedido.Email,
		&pedido.Telefone,
		&pedido.Endereco,
		&pedido.Numero,
		&pedido.Complemento,
		&pedido.Bairro,
		&pedido.Cidade,
		&pedido.Estado,
		&pedido.CEP,
		&pedido.FormaPagamento,
		&pedido.Status,
		&pedido.Total,
		&pedido.CriadoEm,
		&pedido.AtualizadoEm,
	)

	if erro != nil {
		return modelos.Pedido{}, erro
	}

	return pedido, nil
}

// BuscarItensPorPedidoID busca os itens de um pedido
func (repositorio Pedidos) BuscarItensPorPedidoID(pedidoID uint64) ([]modelos.ItemPedido, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT 
			id, pedido_id, produto_id, nome_produto, preco_unitario,
			quantidade, tamanho, subtotal, foto_url, criadoEm
		FROM itens_pedido
		WHERE pedido_id = ?
		ORDER BY id
	`, pedidoID)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	itens := make([]modelos.ItemPedido, 0)

	for linhas.Next() {
		var item modelos.ItemPedido
		if erro = linhas.Scan(
			&item.ID,
			&item.PedidoID,
			&item.ProdutoID,
			&item.NomeProduto,
			&item.PrecoUnitario,
			&item.Quantidade,
			&item.Tamanho,
			&item.Subtotal,
			&item.FotoURL,
			&item.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		itens = append(itens, item)
	}

	return itens, nil
}

// BuscarPedidoCompleto busca um pedido completo com seus itens
func (repositorio Pedidos) BuscarPedidoCompleto(pedidoID uint64) (modelos.PedidoCompleto, error) {
	pedido, erro := repositorio.BuscarPorID(pedidoID)
	if erro != nil {
		return modelos.PedidoCompleto{}, erro
	}

	itens, erro := repositorio.BuscarItensPorPedidoID(pedidoID)
	if erro != nil {
		return modelos.PedidoCompleto{}, erro
	}

	return modelos.PedidoCompleto{
		Pedido: pedido,
		Itens:  itens,
	}, nil
}

// BuscarPorUsuario busca todos os pedidos de um usuário
func (repositorio Pedidos) BuscarPorUsuario(usuarioID uint64) ([]modelos.Pedido, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT 
			id, usuario_id, nome_completo, email, telefone,
			endereco, numero, complemento, bairro, cidade, estado, cep,
			forma_pagamento, status, total, criadoEm, atualizadoEm
		FROM pedidos
		WHERE usuario_id = ?
		ORDER BY criadoEm DESC
	`, usuarioID)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	pedidos := make([]modelos.Pedido, 0)

	for linhas.Next() {
		var pedido modelos.Pedido
		if erro = linhas.Scan(
			&pedido.ID,
			&pedido.UsuarioID,
			&pedido.NomeCompleto,
			&pedido.Email,
			&pedido.Telefone,
			&pedido.Endereco,
			&pedido.Numero,
			&pedido.Complemento,
			&pedido.Bairro,
			&pedido.Cidade,
			&pedido.Estado,
			&pedido.CEP,
			&pedido.FormaPagamento,
			&pedido.Status,
			&pedido.Total,
			&pedido.CriadoEm,
			&pedido.AtualizadoEm,
		); erro != nil {
			return nil, erro
		}

		pedidos = append(pedidos, pedido)
	}

	return pedidos, nil
}

// AtualizarStatus atualiza o status de um pedido
func (repositorio Pedidos) AtualizarStatus(pedidoID uint64, status string) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE pedidos SET status = ? WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(status, pedidoID)
	if erro != nil {
		return erro
	}

	linhasAfetadas, erro := resultado.RowsAffected()
	if erro != nil {
		return erro
	}

	if linhasAfetadas == 0 {
		return fmt.Errorf("pedido não encontrado")
	}

	return nil
}

// ListarTodos lista todos os pedidos (admin)
func (repositorio Pedidos) ListarTodos() ([]modelos.Pedido, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT 
			id, usuario_id, nome_completo, email, telefone,
			endereco, numero, complemento, bairro, cidade, estado, cep,
			forma_pagamento, status, total, criadoEm, atualizadoEm
		FROM pedidos
		ORDER BY criadoEm DESC
	`)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	pedidos := make([]modelos.Pedido, 0)

	for linhas.Next() {
		var pedido modelos.Pedido
		if erro = linhas.Scan(
			&pedido.ID,
			&pedido.UsuarioID,
			&pedido.NomeCompleto,
			&pedido.Email,
			&pedido.Telefone,
			&pedido.Endereco,
			&pedido.Numero,
			&pedido.Complemento,
			&pedido.Bairro,
			&pedido.Cidade,
			&pedido.Estado,
			&pedido.CEP,
			&pedido.FormaPagamento,
			&pedido.Status,
			&pedido.Total,
			&pedido.CriadoEm,
			&pedido.AtualizadoEm,
		); erro != nil {
			return nil, erro
		}

		pedidos = append(pedidos, pedido)
	}

	return pedidos, nil
}
