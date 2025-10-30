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

	// Buscar itens do carrinho do usuário e armazenar em slice
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

	// Armazenar os itens em um slice antes de fechar a query
	type ItemTemp struct {
		ProdutoID  uint64
		Quantidade int
		Nome       string
		Preco      float64
		Tamanho    string
		FotoURL    string
	}

	var itensTemp []ItemTemp

	for linhas.Next() {
		var item ItemTemp
		if erro = linhas.Scan(&item.ProdutoID, &item.Quantidade, &item.Nome, &item.Preco, &item.Tamanho, &item.FotoURL); erro != nil {
			linhas.Close()
			return 0, erro
		}
		itensTemp = append(itensTemp, item)
	}
	linhas.Close()

	// Verificar se há itens no carrinho
	if len(itensTemp) == 0 {
		return 0, fmt.Errorf("carrinho vazio - não é possível criar pedido sem itens")
	}

	// Inserir os itens do pedido
	for _, item := range itensTemp {
		subtotal := item.Preco * float64(item.Quantidade)

		_, erro = tx.Exec(`
			INSERT INTO itens_pedido (
				pedido_id, produto_id, nome_produto, preco_unitario,
				quantidade, tamanho, subtotal, foto_url
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, pedidoID, item.ProdutoID, item.Nome, item.Preco, item.Quantidade, item.Tamanho, subtotal, item.FotoURL)

		if erro != nil {
			return 0, erro
		}
	}

	// Limpar carrinho do usuário se houver itens
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
			forma_pagamento, status, codigo_rastreio, total, criadoEm, atualizadoEm
		FROM pedidos
		WHERE id = ?
	`, pedidoID)

	var pedido modelos.Pedido
	var codigoRastreio sql.NullString

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
		&codigoRastreio,
		&pedido.Total,
		&pedido.CriadoEm,
		&pedido.AtualizadoEm,
	)

	if erro != nil {
		return modelos.Pedido{}, erro
	}

	if codigoRastreio.Valid {
		pedido.CodigoRastreio = codigoRastreio.String
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
			forma_pagamento, status, codigo_rastreio, total, criadoEm, atualizadoEm
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
		var codigoRastreio sql.NullString

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
			&codigoRastreio,
			&pedido.Total,
			&pedido.CriadoEm,
			&pedido.AtualizadoEm,
		); erro != nil {
			return nil, erro
		}

		if codigoRastreio.Valid {
			pedido.CodigoRastreio = codigoRastreio.String
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
			forma_pagamento, status, codigo_rastreio, total, criadoEm, atualizadoEm
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
		var codigoRastreio sql.NullString

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
			&codigoRastreio,
			&pedido.Total,
			&pedido.CriadoEm,
			&pedido.AtualizadoEm,
		); erro != nil {
			return nil, erro
		}

		if codigoRastreio.Valid {
			pedido.CodigoRastreio = codigoRastreio.String
		}

		pedidos = append(pedidos, pedido)
	}

	return pedidos, nil
}

// AtualizarCodigoRastreio atualiza o código de rastreio de um pedido e muda status para "enviado"
func (repositorio Pedidos) AtualizarCodigoRastreio(pedidoID uint64, codigoRastreio string) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE pedidos SET codigo_rastreio = ?, status = 'enviado' WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(codigoRastreio, pedidoID)
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

// ConfirmarEntrega marca o pedido como recebido
func (repositorio Pedidos) ConfirmarEntrega(pedidoID uint64, usuarioID uint64) error {
	// Verificar se o pedido pertence ao usuário
	pedido, erro := repositorio.BuscarPorID(pedidoID)
	if erro != nil {
		return erro
	}

	if pedido.UsuarioID != usuarioID {
		return fmt.Errorf("pedido não pertence ao usuário")
	}

	statement, erro := repositorio.db.Prepare(
		"UPDATE pedidos SET status = 'recebido' WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(pedidoID)
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
