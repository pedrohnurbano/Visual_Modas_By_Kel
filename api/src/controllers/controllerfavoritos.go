package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// AdicionarFavorito adiciona um produto aos favoritos do usuário
func AdicionarFavorito(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do usuário do token
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Ler corpo da requisição
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	var favorito modelos.Favorito
	if erro = json.Unmarshal(corpoRequisicao, &favorito); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Conectar ao banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Verificar se o produto existe
	repositorioProdutos := repositorios.NovoRepositorioDeProdutos(db)
	_, erro = repositorioProdutos.BuscarPorID(favorito.ProdutoID)
	if erro != nil {
		respostas.Erro(w, http.StatusNotFound, erro)
		return
	}

	// Adicionar favorito
	repositorio := repositorios.NovoRepositorioDeFavoritos(db)
	if erro = repositorio.Adicionar(usuarioID, favorito.ProdutoID); erro != nil {
		// Log do erro para debug
		println("ERRO ao adicionar favorito:", erro.Error())
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Log de sucesso
	println("Favorito adicionado com sucesso! UsuarioID:", usuarioID, "ProdutoID:", favorito.ProdutoID)
	respostas.JSON(w, http.StatusCreated, map[string]string{"mensagem": "Produto adicionado aos favoritos"})
}

// RemoverFavorito remove um produto dos favoritos do usuário
func RemoverFavorito(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do usuário do token
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Extrair ID do produto dos parâmetros da URL
	parametros := mux.Vars(r)
	produtoID, erro := strconv.ParseUint(parametros["produtoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Conectar ao banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Remover favorito
	repositorio := repositorios.NovoRepositorioDeFavoritos(db)
	if erro = repositorio.Remover(usuarioID, produtoID); erro != nil {
		println("ERRO ao remover favorito:", erro.Error())
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	println("Favorito removido com sucesso! UsuarioID:", usuarioID, "ProdutoID:", produtoID)
	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Produto removido dos favoritos"})
}

// BuscarFavoritos retorna todos os favoritos do usuário
func BuscarFavoritos(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do usuário do token
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Conectar ao banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Buscar favoritos
	repositorio := repositorios.NovoRepositorioDeFavoritos(db)
	favoritos, erro := repositorio.BuscarPorUsuario(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, favoritos)
}

// BuscarIDsFavoritos retorna apenas os IDs dos produtos favoritos do usuário
func BuscarIDsFavoritos(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do usuário do token
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Conectar ao banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Buscar IDs dos favoritos
	repositorio := repositorios.NovoRepositorioDeFavoritos(db)
	produtoIDs, erro := repositorio.BuscarIDsFavoritosPorUsuario(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, produtoIDs)
}

// VerificarFavorito verifica se um produto está nos favoritos do usuário
func VerificarFavorito(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do usuário do token
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Extrair ID do produto dos parâmetros da URL
	parametros := mux.Vars(r)
	produtoID, erro := strconv.ParseUint(parametros["produtoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Conectar ao banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Verificar se é favorito
	repositorio := repositorios.NovoRepositorioDeFavoritos(db)
	isFavorito, erro := repositorio.VerificarFavorito(usuarioID, produtoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, map[string]bool{"isFavorito": isFavorito})
}
