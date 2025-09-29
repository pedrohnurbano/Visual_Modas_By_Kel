package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// CriarProduto insere um novo produto no banco de dados
func CriarProduto(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var produto modelos.Produto
	if erro = json.Unmarshal(corpoRequest, &produto); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	produto.UsuarioID = usuarioID

	// Processar imagem base64
	if strings.HasPrefix(produto.FotoURL, "data:image") {
		nomeArquivo, erro := salvarImagemBase64(produto.FotoURL)
		if erro != nil {
			respostas.Erro(w, http.StatusInternalServerError, 
				errors.New("erro ao salvar imagem: " + erro.Error()))
			return
		}
		produto.FotoURL = nomeArquivo
	}

	if erro = produto.Preparar("cadastro"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produto.ID, erro = repositorio.Criar(produto)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusCreated, produto)
}

// salvarImagemBase64 decodifica base64 e salva como arquivo
func salvarImagemBase64(dataURL string) (string, error) {
	// Separar o prefixo data:image/xxx;base64, dos dados
	parts := strings.Split(dataURL, ",")
	if len(parts) != 2 {
		return "", errors.New("formato de imagem base64 inválido")
	}

	// Extrair tipo de imagem (png, jpg, etc)
	tipoImagem := "jpg"
	if strings.Contains(parts[0], "image/png") {
		tipoImagem = "png"
	} else if strings.Contains(parts[0], "image/webp") {
		tipoImagem = "webp"
	}

	// Decodificar base64
	imagemDecodificada, erro := base64.StdEncoding.DecodeString(parts[1])
	if erro != nil {
		return "", erro
	}

	// Criar pasta uploads se não existir
	pastaUploads := "./uploads/produtos"
	if erro := os.MkdirAll(pastaUploads, 0755); erro != nil {
		return "", erro
	}

	// Gerar nome único para o arquivo
	nomeArquivo := fmt.Sprintf("produto_%d.%s", time.Now().Unix(), tipoImagem)
	caminhoCompleto := filepath.Join(pastaUploads, nomeArquivo)

	// Salvar arquivo
	arquivo, erro := os.Create(caminhoCompleto)
	if erro != nil {
		return "", erro
	}
	defer arquivo.Close()

	if _, erro := arquivo.Write(imagemDecodificada); erro != nil {
		return "", erro
	}

	// Retornar apenas o caminho relativo
	return fmt.Sprintf("/uploads/produtos/%s", nomeArquivo), nil
}

// BuscarProdutos retorna todos os produtos ou filtrados
func BuscarProdutos(w http.ResponseWriter, r *http.Request) {
	filtro := strings.ToLower(r.URL.Query().Get("filtro"))
	
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produtos, erro := repositorio.Buscar(filtro)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, produtos)
}

// BuscarProduto retorna um produto específico
func BuscarProduto(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	produtoID, erro := strconv.ParseUint(parametros["produtoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produto, erro := repositorio.BuscarPorID(produtoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, produto)
}

// BuscarProdutosPorCategoria retorna produtos de uma categoria
func BuscarProdutosPorCategoria(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	categoria := parametros["categoria"]

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produtos, erro := repositorio.BuscarPorCategoria(categoria)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, produtos)
}

// AtualizarProduto modifica as informações de um produto
func AtualizarProduto(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	produtoID, erro := strconv.ParseUint(parametros["produtoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produtoNoBanco, erro := repositorio.BuscarPorID(produtoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	repositorioUsuarios := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorioUsuarios.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if produtoNoBanco.UsuarioID != usuarioID && usuario.Role != "admin" {
		respostas.Erro(w, http.StatusForbidden, errors.New("você não tem permissão para editar este produto"))
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var produto modelos.Produto
	if erro = json.Unmarshal(corpoRequisicao, &produto); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Processar imagem base64 se houver
	if strings.HasPrefix(produto.FotoURL, "data:image") {
		// Deletar imagem antiga se existir
		if produtoNoBanco.FotoURL != "" && !strings.HasPrefix(produtoNoBanco.FotoURL, "http") {
			os.Remove("." + produtoNoBanco.FotoURL)
		}
		
		nomeArquivo, erro := salvarImagemBase64(produto.FotoURL)
		if erro != nil {
			respostas.Erro(w, http.StatusInternalServerError, 
				errors.New("erro ao salvar imagem: " + erro.Error()))
			return
		}
		produto.FotoURL = nomeArquivo
	}

	if erro = produto.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.Atualizar(produtoID, produto); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarProduto remove um produto (soft delete)
func DeletarProduto(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	produtoID, erro := strconv.ParseUint(parametros["produtoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produtoNoBanco, erro := repositorio.BuscarPorID(produtoID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	repositorioUsuarios := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorioUsuarios.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if produtoNoBanco.UsuarioID != usuarioID && usuario.Role != "admin" {
		respostas.Erro(w, http.StatusForbidden, errors.New("você não tem permissão para deletar este produto"))
		return
	}

	if erro = repositorio.Deletar(produtoID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// MeusProdutos retorna os produtos do usuário autenticado
func MeusProdutos(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeProdutos(db)
	produtos, erro := repositorio.BuscarPorUsuario(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, produtos)
}
