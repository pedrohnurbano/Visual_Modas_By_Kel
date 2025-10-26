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

// AdicionarAoCarrinho adiciona um produto ao carrinho do usuário
func AdicionarAoCarrinho(w http.ResponseWriter, r *http.Request) {
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

	var item modelos.ItemCarrinho
	if erro = json.Unmarshal(corpoRequisicao, &item); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Validar quantidade
	if item.Quantidade <= 0 {
		item.Quantidade = 1
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
	_, erro = repositorioProdutos.BuscarPorID(item.ProdutoID)
	if erro != nil {
		respostas.Erro(w, http.StatusNotFound, erro)
		return
	}

	// Adicionar ao carrinho
	repositorio := repositorios.NovoRepositorioDeCarrinho(db)
	if erro = repositorio.Adicionar(usuarioID, item.ProdutoID, item.Quantidade); erro != nil {
		println("ERRO ao adicionar ao carrinho:", erro.Error())
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	println("Produto adicionado ao carrinho! UsuarioID:", usuarioID, "ProdutoID:", item.ProdutoID, "Quantidade:", item.Quantidade)
	respostas.JSON(w, http.StatusCreated, map[string]string{"mensagem": "Produto adicionado ao carrinho"})
}

// AtualizarQuantidadeCarrinho atualiza a quantidade de um item no carrinho
func AtualizarQuantidadeCarrinho(w http.ResponseWriter, r *http.Request) {
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

	// Ler corpo da requisição
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	var dados struct {
		Quantidade int `json:"quantidade"`
	}
	if erro = json.Unmarshal(corpoRequisicao, &dados); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Validar quantidade
	if dados.Quantidade <= 0 {
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

	// Atualizar quantidade
	repositorio := repositorios.NovoRepositorioDeCarrinho(db)
	if erro = repositorio.AtualizarQuantidade(usuarioID, produtoID, dados.Quantidade); erro != nil {
		println("ERRO ao atualizar quantidade no carrinho:", erro.Error())
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	println("Quantidade atualizada no carrinho! UsuarioID:", usuarioID, "ProdutoID:", produtoID, "Quantidade:", dados.Quantidade)
	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Quantidade atualizada"})
}

// RemoverDoCarrinho remove um produto do carrinho do usuário
func RemoverDoCarrinho(w http.ResponseWriter, r *http.Request) {
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

	// Remover do carrinho
	repositorio := repositorios.NovoRepositorioDeCarrinho(db)
	if erro = repositorio.Remover(usuarioID, produtoID); erro != nil {
		println("ERRO ao remover do carrinho:", erro.Error())
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	println("Produto removido do carrinho! UsuarioID:", usuarioID, "ProdutoID:", produtoID)
	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Produto removido do carrinho"})
}

// LimparCarrinho remove todos os itens do carrinho do usuário
func LimparCarrinho(w http.ResponseWriter, r *http.Request) {
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

	// Limpar carrinho
	repositorio := repositorios.NovoRepositorioDeCarrinho(db)
	if erro = repositorio.LimparCarrinho(usuarioID); erro != nil {
		println("ERRO ao limpar carrinho:", erro.Error())
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	println("Carrinho limpo! UsuarioID:", usuarioID)
	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Carrinho limpo"})
}

// BuscarCarrinho retorna todos os itens do carrinho do usuário
func BuscarCarrinho(w http.ResponseWriter, r *http.Request) {
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

	// Buscar itens do carrinho
	repositorio := repositorios.NovoRepositorioDeCarrinho(db)
	itens, erro := repositorio.BuscarPorUsuario(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, itens)
}

// BuscarResumoCarrinho retorna o resumo do carrinho com totais
func BuscarResumoCarrinho(w http.ResponseWriter, r *http.Request) {
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

	// Buscar resumo do carrinho
	repositorio := repositorios.NovoRepositorioDeCarrinho(db)
	resumo, erro := repositorio.BuscarResumoCarrinho(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, resumo)
}

// ContarItensCarrinho retorna a quantidade total de itens no carrinho
func ContarItensCarrinho(w http.ResponseWriter, r *http.Request) {
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

	// Contar itens
	repositorio := repositorios.NovoRepositorioDeCarrinho(db)
	total, erro := repositorio.ContarItens(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, map[string]int{"total": total})
}
