package repositorios

import (
	"api/src/modelos"
	"database/sql"
)

// Favoritos representa um repositório de favoritos
type Favoritos struct {
	db *sql.DB
}

// NovoRepositorioDeFavoritos cria um repositório de favoritos
func NovoRepositorioDeFavoritos(db *sql.DB) *Favoritos {
	return &Favoritos{db}
}

// Adicionar adiciona um produto aos favoritos do usuário
func (repositorio Favoritos) Adicionar(usuarioID, produtoID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO favoritos (usuario_id, produto_id) VALUES (?, ?)",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(usuarioID, produtoID)
	return erro
}

// Remover remove um produto dos favoritos do usuário
func (repositorio Favoritos) Remover(usuarioID, produtoID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"DELETE FROM favoritos WHERE usuario_id = ? AND produto_id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(usuarioID, produtoID)
	return erro
}

// BuscarPorUsuario retorna todos os favoritos de um usuário com os dados dos produtos
func (repositorio Favoritos) BuscarPorUsuario(usuarioID uint64) ([]modelos.FavoritoComProduto, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT 
			f.id, f.usuario_id, f.produto_id, f.criadoEm,
			p.id, p.nome, p.descricao, p.preco, p.tamanho, p.categoria, p.secao, 
			p.foto_url, p.usuario_id, p.ativo, p.criadoEm, p.atualizadoEm
		FROM favoritos f
		INNER JOIN produtos p ON f.produto_id = p.id
		WHERE f.usuario_id = ? AND p.ativo = true
		ORDER BY f.criadoEm DESC
	`, usuarioID)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	// Inicializar com slice vazio em vez de nil
	favoritos := make([]modelos.FavoritoComProduto, 0)

	for linhas.Next() {
		var favorito modelos.FavoritoComProduto
		if erro = linhas.Scan(
			&favorito.ID,
			&favorito.UsuarioID,
			&favorito.ProdutoID,
			&favorito.CriadoEm,
			&favorito.Produto.ID,
			&favorito.Produto.Nome,
			&favorito.Produto.Descricao,
			&favorito.Produto.Preco,
			&favorito.Produto.Tamanho,
			&favorito.Produto.Categoria,
			&favorito.Produto.Secao,
			&favorito.Produto.FotoURL,
			&favorito.Produto.UsuarioID,
			&favorito.Produto.Ativo,
			&favorito.Produto.CriadoEm,
			&favorito.Produto.AtualizadoEm,
		); erro != nil {
			return nil, erro
		}

		favoritos = append(favoritos, favorito)
	}

	return favoritos, nil
}

// VerificarFavorito verifica se um produto está nos favoritos do usuário
func (repositorio Favoritos) VerificarFavorito(usuarioID, produtoID uint64) (bool, error) {
	var count int
	erro := repositorio.db.QueryRow(
		"SELECT COUNT(*) FROM favoritos WHERE usuario_id = ? AND produto_id = ?",
		usuarioID, produtoID,
	).Scan(&count)

	if erro != nil {
		return false, erro
	}

	return count > 0, nil
}

// BuscarIDsFavoritosPorUsuario retorna apenas os IDs dos produtos favoritos do usuário
func (repositorio Favoritos) BuscarIDsFavoritosPorUsuario(usuarioID uint64) ([]uint64, error) {
	linhas, erro := repositorio.db.Query(`
		SELECT produto_id FROM favoritos WHERE usuario_id = ?
	`, usuarioID)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	// Inicializar com slice vazio em vez de nil
	produtoIDs := make([]uint64, 0)

	for linhas.Next() {
		var produtoID uint64
		if erro = linhas.Scan(&produtoID); erro != nil {
			return nil, erro
		}
		produtoIDs = append(produtoIDs, produtoID)
	}

	return produtoIDs, nil
}
