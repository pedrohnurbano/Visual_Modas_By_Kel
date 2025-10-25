package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
)

// Produtos representa um repositório de produtos
type Produtos struct {
	db *sql.DB
}

// NovoRepositorioDeProdutos cria um novo repositório de produtos
func NovoRepositorioDeProdutos(db *sql.DB) *Produtos {
	return &Produtos{db}
}

// Criar insere um produto no banco de dados
func (repositorio Produtos) Criar(produto modelos.Produto) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO produtos (nome, descricao, preco, tamanho, categoria, secao, genero, foto_url, usuario_id, ativo) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(
		produto.Nome,
		produto.Descricao,
		produto.Preco,
		produto.Tamanho,
		produto.Categoria,
		produto.Secao,
		produto.Genero,
		produto.FotoURL,
		produto.UsuarioID,
		true, // ativo por padrão
	)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

// Buscar retorna todos os produtos que atendem um filtro
func (repositorio Produtos) Buscar(filtro string) ([]modelos.Produto, error) {
	filtro = fmt.Sprintf("%%%s%%", filtro)

	linhas, erro := repositorio.db.Query(
		`SELECT id, nome, descricao, preco, tamanho, categoria, secao, genero, foto_url, usuario_id, ativo, criadoEm, atualizadoEm 
		FROM produtos 
		WHERE ativo = true AND (nome LIKE ? OR descricao LIKE ? OR categoria LIKE ? OR secao LIKE ? OR genero LIKE ?)
		ORDER BY criadoEm DESC`,
		filtro, filtro, filtro, filtro, filtro,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var produtos []modelos.Produto

	for linhas.Next() {
		var produto modelos.Produto
		if erro = linhas.Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Descricao,
			&produto.Preco,
			&produto.Tamanho,
			&produto.Categoria,
			&produto.Secao,
			&produto.Genero,
			&produto.FotoURL,
			&produto.UsuarioID,
			&produto.Ativo,
			&produto.CriadoEm,
			&produto.AtualizadoEm,
		); erro != nil {
			return nil, erro
		}

		produtos = append(produtos, produto)
	}

	return produtos, nil
}

// BuscarPorID retorna um produto específico
func (repositorio Produtos) BuscarPorID(ID uint64) (modelos.Produto, error) {
	linha := repositorio.db.QueryRow(
		`SELECT id, nome, descricao, preco, tamanho, categoria, secao, genero, foto_url, usuario_id, ativo, criadoEm, atualizadoEm 
		FROM produtos WHERE id = ?`,
		ID,
	)

	var produto modelos.Produto
	erro := linha.Scan(
		&produto.ID,
		&produto.Nome,
		&produto.Descricao,
		&produto.Preco,
		&produto.Tamanho,
		&produto.Categoria,
		&produto.Secao,
		&produto.Genero,
		&produto.FotoURL,
		&produto.UsuarioID,
		&produto.Ativo,
		&produto.CriadoEm,
		&produto.AtualizadoEm,
	)

	return produto, erro
}

// BuscarPorCategoria retorna produtos de uma categoria específica
func (repositorio Produtos) BuscarPorCategoria(categoria string) ([]modelos.Produto, error) {
	linhas, erro := repositorio.db.Query(
		`SELECT id, nome, descricao, preco, tamanho, categoria, secao, genero, foto_url, usuario_id, ativo, criadoEm, atualizadoEm 
		FROM produtos 
		WHERE ativo = true AND categoria = ?
		ORDER BY criadoEm DESC`,
		categoria,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var produtos []modelos.Produto

	for linhas.Next() {
		var produto modelos.Produto
		if erro = linhas.Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Descricao,
			&produto.Preco,
			&produto.Tamanho,
			&produto.Categoria,
			&produto.Secao,
			&produto.Genero,
			&produto.FotoURL,
			&produto.UsuarioID,
			&produto.Ativo,
			&produto.CriadoEm,
			&produto.AtualizadoEm,
		); erro != nil {
			return nil, erro
		}

		produtos = append(produtos, produto)
	}

	return produtos, nil
}

// BuscarPorSecao retorna produtos de uma seção específica
func (repositorio Produtos) BuscarPorSecao(secao string) ([]modelos.Produto, error) {
	linhas, erro := repositorio.db.Query(
		`SELECT id, nome, descricao, preco, tamanho, categoria, secao, genero, foto_url, usuario_id, ativo, criadoEm, atualizadoEm 
		FROM produtos 
		WHERE ativo = true AND secao = ?
		ORDER BY criadoEm DESC`,
		secao,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var produtos []modelos.Produto

	for linhas.Next() {
		var produto modelos.Produto
		if erro = linhas.Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Descricao,
			&produto.Preco,
			&produto.Tamanho,
			&produto.Categoria,
			&produto.Secao,
			&produto.Genero,
			&produto.FotoURL,
			&produto.UsuarioID,
			&produto.Ativo,
			&produto.CriadoEm,
			&produto.AtualizadoEm,
		); erro != nil {
			return nil, erro
		}

		produtos = append(produtos, produto)
	}

	return produtos, nil
}

// BuscarPorUsuario retorna todos os produtos de um usuário
func (repositorio Produtos) BuscarPorUsuario(usuarioID uint64) ([]modelos.Produto, error) {
	linhas, erro := repositorio.db.Query(
		`SELECT id, nome, descricao, preco, tamanho, categoria, secao, genero, foto_url, usuario_id, ativo, criadoEm, atualizadoEm 
		FROM produtos 
		WHERE usuario_id = ?
		ORDER BY criadoEm DESC`,
		usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var produtos []modelos.Produto

	for linhas.Next() {
		var produto modelos.Produto
		if erro = linhas.Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Descricao,
			&produto.Preco,
			&produto.Tamanho,
			&produto.Categoria,
			&produto.Secao,
			&produto.Genero,
			&produto.FotoURL,
			&produto.UsuarioID,
			&produto.Ativo,
			&produto.CriadoEm,
			&produto.AtualizadoEm,
		); erro != nil {
			return nil, erro
		}

		produtos = append(produtos, produto)
	}

	return produtos, nil
}

// Atualizar modifica as informações de um produto
func (repositorio Produtos) Atualizar(ID uint64, produto modelos.Produto) error {
	statement, erro := repositorio.db.Prepare(
		`UPDATE produtos 
		SET nome = ?, descricao = ?, preco = ?, tamanho = ?, categoria = ?, secao = ?, genero = ?, foto_url = ? 
		WHERE id = ?`,
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(
		produto.Nome,
		produto.Descricao,
		produto.Preco,
		produto.Tamanho,
		produto.Categoria,
		produto.Secao,
		produto.Genero,
		produto.FotoURL,
		ID,
	); erro != nil {
		return erro
	}

	return nil
}

// Deletar remove um produto do banco (soft delete)
func (repositorio Produtos) Deletar(ID uint64) error {
	statement, erro := repositorio.db.Prepare("UPDATE produtos SET ativo = false WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}

// DeletarPermanente remove permanentemente um produto do banco
func (repositorio Produtos) DeletarPermanente(ID uint64) error {
	statement, erro := repositorio.db.Prepare("DELETE FROM produtos WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}

// BuscarComFiltros retorna produtos filtrados por categoria, tamanho, gênero e termo de busca
func (repositorio Produtos) BuscarComFiltros(categoria, tamanho, genero, termoBusca string) ([]modelos.Produto, error) {
	query := `SELECT id, nome, descricao, preco, tamanho, categoria, secao, genero, foto_url, usuario_id, ativo, criadoEm, atualizadoEm 
		FROM produtos 
		WHERE ativo = true`

	var args []interface{}

	// Adicionar filtro de categoria
	if categoria != "" && categoria != "all" {
		query += " AND categoria = ?"
		args = append(args, categoria)
	}

	// Adicionar filtro de tamanho
	if tamanho != "" && tamanho != "all" {
		query += " AND tamanho = ?"
		args = append(args, tamanho)
	}

	// Adicionar filtro de gênero
	if genero != "" && genero != "all" {
		query += " AND genero = ?"
		args = append(args, genero)
	}

	// Adicionar termo de busca
	if termoBusca != "" {
		termoBusca = fmt.Sprintf("%%%s%%", termoBusca)
		query += " AND (nome LIKE ? OR descricao LIKE ?)"
		args = append(args, termoBusca, termoBusca)
	}

	query += " ORDER BY criadoEm DESC"

	linhas, erro := repositorio.db.Query(query, args...)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var produtos []modelos.Produto

	for linhas.Next() {
		var produto modelos.Produto
		if erro = linhas.Scan(
			&produto.ID,
			&produto.Nome,
			&produto.Descricao,
			&produto.Preco,
			&produto.Tamanho,
			&produto.Categoria,
			&produto.Secao,
			&produto.Genero,
			&produto.FotoURL,
			&produto.UsuarioID,
			&produto.Ativo,
			&produto.CriadoEm,
			&produto.AtualizadoEm,
		); erro != nil {
			return nil, erro
		}

		produtos = append(produtos, produto)
	}

	return produtos, nil
}
