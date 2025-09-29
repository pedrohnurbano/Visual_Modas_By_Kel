package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/config"
	"visual_modas_by_kel/visual_modas_by_kel/src/cookies"
	"visual_modas_by_kel/visual_modas_by_kel/src/modelos"
	"visual_modas_by_kel/visual_modas_by_kel/src/requisicoes"
	"visual_modas_by_kel/visual_modas_by_kel/src/respostas"
	"visual_modas_by_kel/visual_modas_by_kel/src/utils"

	"github.com/gorilla/mux"
)

// CarregarPaginaCadastroProduto carrega a página de cadastro de produto
func CarregarPaginaCadastroProduto(w http.ResponseWriter, r *http.Request) {
	_, erro := cookies.Ler(r)
	if erro != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	utils.ExecutarTemplate(w, "cadastro-produto.html", nil)
}

// CarregarPaginaProdutos carrega a página com lista de produtos
func CarregarPaginaProdutos(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "produtos.html", nil)
}

// CarregarPaginaMeusProdutos carrega a página com produtos do usuário
func CarregarPaginaMeusProdutos(w http.ResponseWriter, r *http.Request) {
	_, erro := cookies.Ler(r)
	if erro != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	utils.ExecutarTemplate(w, "meus-produtos.html", nil)
}

// CriarProduto chama a API para cadastrar um produto
func CriarProduto(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Montar JSON do produto
	produto := map[string]interface{}{
		"nome":      r.FormValue("nome"),
		"descricao": r.FormValue("descricao"),
		"preco":     r.FormValue("preco"),
		"tamanho":   r.FormValue("tamanho"),
		"categoria": r.FormValue("categoria"),
		"foto_url":  r.FormValue("foto_url"), // Já vem em base64 do frontend
	}

	produtoJSON, erro := json.Marshal(produto)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/produtos", config.APIURL)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPost, url, produtoJSON)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	respostas.JSON(w, response.StatusCode, nil)
}

// BuscarProdutos busca todos os produtos
func BuscarProdutos(w http.ResponseWriter, r *http.Request) {
	filtro := r.URL.Query().Get("filtro")
	url := fmt.Sprintf("%s/produtos?filtro=%s", config.APIURL, filtro)

	response, erro := http.Get(url)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	var produtos []modelos.Produto
	if erro = json.NewDecoder(response.Body).Decode(&produtos); erro != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	respostas.JSON(w, http.StatusOK, produtos)
}

// BuscarProduto busca um produto específico
func BuscarProduto(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	produtoID := parametros["produtoId"]

	url := fmt.Sprintf("%s/produtos/%s", config.APIURL, produtoID)
	response, erro := http.Get(url)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	var produto modelos.Produto
	if erro = json.NewDecoder(response.Body).Decode(&produto); erro != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	respostas.JSON(w, http.StatusOK, produto)
}

// BuscarMeusProdutos busca os produtos do usuário logado
func BuscarMeusProdutos(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/meus-produtos", config.APIURL)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	var produtos []modelos.Produto
	if erro = json.NewDecoder(response.Body).Decode(&produtos); erro != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	respostas.JSON(w, http.StatusOK, produtos)
}

// AtualizarProduto atualiza um produto
func AtualizarProduto(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	produtoID := parametros["produtoId"]

	r.ParseForm()

	produto := map[string]interface{}{
		"nome":      r.FormValue("nome"),
		"descricao": r.FormValue("descricao"),
		"preco":     r.FormValue("preco"),
		"tamanho":   r.FormValue("tamanho"),
		"categoria": r.FormValue("categoria"),
		"foto_url":  r.FormValue("foto_url"),
	}

	produtoJSON, erro := json.Marshal(produto)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/produtos/%s", config.APIURL, produtoID)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPut, url, produtoJSON)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	respostas.JSON(w, response.StatusCode, nil)
}

// DeletarProduto deleta um produto
func DeletarProduto(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	produtoID := parametros["produtoId"]

	url := fmt.Sprintf("%s/produtos/%s", config.APIURL, produtoID)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodDelete, url, nil)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	respostas.JSON(w, response.StatusCode, nil)
}
