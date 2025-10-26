package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/config"
	"visual_modas_by_kel/visual_modas_by_kel/src/requisicoes"
	"visual_modas_by_kel/visual_modas_by_kel/src/respostas"

	"github.com/gorilla/mux"
)

// AdicionarAoCarrinho chama a API para adicionar um produto ao carrinho
func AdicionarAoCarrinho(w http.ResponseWriter, r *http.Request) {
	body, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: "Erro ao ler dados"})
		return
	}
	defer r.Body.Close()

	url := fmt.Sprintf("%s/carrinho", config.APIURL)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPost, url, body)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	var resultado map[string]interface{}
	json.NewDecoder(response.Body).Decode(&resultado)
	respostas.JSON(w, http.StatusCreated, resultado)
}

// AtualizarQuantidadeCarrinho chama a API para atualizar quantidade de um item
func AtualizarQuantidadeCarrinho(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	produtoID := parametros["produtoId"]

	body, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: "Erro ao ler dados"})
		return
	}
	defer r.Body.Close()

	url := fmt.Sprintf("%s/carrinho/%s", config.APIURL, produtoID)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPut, url, body)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	var resultado map[string]interface{}
	json.NewDecoder(response.Body).Decode(&resultado)
	respostas.JSON(w, http.StatusOK, resultado)
}

// RemoverDoCarrinho chama a API para remover um produto do carrinho
func RemoverDoCarrinho(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	produtoID := parametros["produtoId"]

	url := fmt.Sprintf("%s/carrinho/%s", config.APIURL, produtoID)
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

	var resultado map[string]interface{}
	json.NewDecoder(response.Body).Decode(&resultado)
	respostas.JSON(w, http.StatusOK, resultado)
}

// LimparCarrinho chama a API para limpar o carrinho
func LimparCarrinho(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/carrinho/limpar", config.APIURL)
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

	var resultado map[string]interface{}
	json.NewDecoder(response.Body).Decode(&resultado)
	respostas.JSON(w, http.StatusOK, resultado)
}

// BuscarCarrinho chama a API para buscar todos os itens do carrinho
func BuscarCarrinho(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/carrinho", config.APIURL)
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

	var itens interface{}
	if erro = json.NewDecoder(response.Body).Decode(&itens); erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: "Erro ao processar resposta"})
		return
	}

	respostas.JSON(w, http.StatusOK, itens)
}

// BuscarResumoCarrinho chama a API para buscar o resumo do carrinho
func BuscarResumoCarrinho(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/carrinho/resumo", config.APIURL)
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

	var resumo interface{}
	if erro = json.NewDecoder(response.Body).Decode(&resumo); erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: "Erro ao processar resposta"})
		return
	}

	respostas.JSON(w, http.StatusOK, resumo)
}

// ContarItensCarrinho chama a API para contar itens do carrinho
func ContarItensCarrinho(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/carrinho/contar", config.APIURL)
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

	var total map[string]int
	if erro = json.NewDecoder(response.Body).Decode(&total); erro != nil {
		// Retornar zero em caso de erro ao parsear
		total = map[string]int{"total": 0}
	}

	respostas.JSON(w, http.StatusOK, total)
}
