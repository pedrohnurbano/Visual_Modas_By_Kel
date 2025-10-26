package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

// Carrinho representa um repositório de carrinho
type Carrinho struct {
	db *sql.DB
}

// NovoRepositorioDeCarrinho cria um repositório de carrinho
func NovoRepositorioDeCarrinho(db *sql.DB) *Carrinho {
	return &Carrinho{db}
}

// Adicionar adiciona um produto ao carrinho do usuário
func (repositorio Carrinho) Adicionar(usuarioID, produtoID uint64, quantidade int) error {
	// Verifica se o item já existe no carrinho
	var itemID uint64
	var quantidadeAtual int
	erro := repositorio.db.QueryRow(
		"SELECT id, quantidade FROM carrinho WHERE usuario_id = ? AND produto_id = ?",
		usuarioID, produtoID,
	).Scan(&itemID, &quantidadeAtual)

	if erro == sql.ErrNoRows {
		// Item não existe, inserir novo
		statement, erro := repositorio.db.Prepare(
			"INSERT INTO carrinho (usuario_id, produto_id, quantidade) VALUES (?, ?, ?)",
		)
		if erro != nil {
			return erro
		}
		defer statement.Close()

		_, erro = statement.Exec(usuarioID, produtoID, quantidade)
		return erro
	} else if erro != nil {
		return erro
	}

	// Item já existe, atualizar quantidade
	statement, erro := repositorio.db.Prepare(
		"UPDATE carrinho SET quantidade = quantidade + ? WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(quantidade, itemID)
	return erro
}

// AtualizarQuantidade atualiza a quantidade de um item no carrinho
func (repositorio Carrinho) AtualizarQuantidade(usuarioID, produtoID uint64, quantidade int) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE carrinho SET quantidade = ? WHERE usuario_id = ? AND produto_id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(quantidade, usuarioID, produtoID)
	return erro
}

// Remover remove um produto do carrinho do usuário
func (repositorio Carrinho) Remover(usuarioID, produtoID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"DELETE FROM carrinho WHERE usuario_id = ? AND produto_id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(usuarioID, produtoID)
	return erro
}

// LimparCarrinho remove todos os itens do carrinho do usuário
func (repositorio Carrinho) LimparCarrinho(usuarioID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"DELETE FROM carrinho WHERE usuario_id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(usuarioID)
	return erro
}

// BuscarPorUsuario retorna todos os itens do carrinho de um usuário com os dados dos produtos
func (repositorio Carrinho) BuscarPorUsuario(usuarioID uint64) ([]modelos.ItemCarrinhoComProduto, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT 
			c.id, c.usuario_id, c.produto_id, c.quantidade, c.criadoEm, c.atualizadoEm,
			p.id, p.nome, p.descricao, p.preco, p.tamanho, p.categoria, p.secao, p.genero,
			p.foto_url, p.usuario_id, p.ativo, p.criadoEm, p.atualizadoEm
		FROM carrinho c
		INNER JOIN produtos p ON c.produto_id = p.id
		WHERE c.usuario_id = ? AND p.ativo = true
		ORDER BY c.criadoEm DESC
	`, usuarioID)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	// Inicializar com slice vazio em vez de nil
	itens := make([]modelos.ItemCarrinhoComProduto, 0)

	for linhas.Next() {
		var item modelos.ItemCarrinhoComProduto
		if erro = linhas.Scan(
			&item.ID,
			&item.UsuarioID,
			&item.ProdutoID,
			&item.Quantidade,
			&item.CriadoEm,
			&item.AtualizadoEm,
			&item.Produto.ID,
			&item.Produto.Nome,
			&item.Produto.Descricao,
			&item.Produto.Preco,
			&item.Produto.Tamanho,
			&item.Produto.Categoria,
			&item.Produto.Secao,
			&item.Produto.Genero,
			&item.Produto.FotoURL,
			&item.Produto.UsuarioID,
			&item.Produto.Ativo,
			&item.Produto.CriadoEm,
			&item.Produto.AtualizadoEm,
		); erro != nil {
			return nil, erro
		}

		itens = append(itens, item)
	}

	return itens, nil
}

// BuscarResumoCarrinho retorna o resumo do carrinho com totais
func (repositorio Carrinho) BuscarResumoCarrinho(usuarioID uint64) (modelos.ResumoCarrinho, error) {
	itens, erro := repositorio.BuscarPorUsuario(usuarioID)
	if erro != nil {
		return modelos.ResumoCarrinho{}, erro
	}

	var resumo modelos.ResumoCarrinho
	resumo.Itens = itens
	resumo.QuantidadeTotal = 0
	resumo.ValorTotal = 0

	for _, item := range itens {
		resumo.QuantidadeTotal += item.Quantidade
		resumo.ValorTotal += item.Produto.Preco * float64(item.Quantidade)
	}

	return resumo, nil
}

// ContarItens retorna a quantidade total de itens no carrinho do usuário
func (repositorio Carrinho) ContarItens(usuarioID uint64) (int, error) {
	var total int
	erro := repositorio.db.QueryRow(
		"SELECT COALESCE(SUM(quantidade), 0) FROM carrinho WHERE usuario_id = ?",
		usuarioID,
	).Scan(&total)

	if erro != nil {
		return 0, erro
	}

	return total, nil
}
