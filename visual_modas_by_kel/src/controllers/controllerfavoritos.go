package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"visual_modas_by_kel/visual_modas_by_kel/src/config"
	"visual_modas_by_kel/visual_modas_by_kel/src/requisicoes"
	"visual_modas_by_kel/visual_modas_by_kel/src/respostas"

	"github.com/gorilla/mux"
)

// AdicionarFavorito chama a API para adicionar um produto aos favoritos
func AdicionarFavorito(w http.ResponseWriter, r *http.Request) {
	body, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: "Erro ao ler dados"})
		return
	}
	defer r.Body.Close()

	url := fmt.Sprintf("%s/favoritos", config.APIURL)
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

	respostas.JSON(w, http.StatusCreated, map[string]string{"mensagem": "Produto adicionado aos favoritos"})
}

// RemoverFavorito chama a API para remover um produto dos favoritos
func RemoverFavorito(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	produtoID := parametros["produtoId"]

	url := fmt.Sprintf("%s/favoritos/%s", config.APIURL, produtoID)
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

	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Produto removido dos favoritos"})
}

// BuscarFavoritos chama a API para buscar todos os favoritos do usuário
func BuscarFavoritos(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/favoritos", config.APIURL)
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

	var favoritos []map[string]interface{}
	if erro = json.NewDecoder(response.Body).Decode(&favoritos); erro != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	respostas.JSON(w, http.StatusOK, favoritos)
}

// BuscarIDsFavoritos chama a API para buscar apenas os IDs dos produtos favoritos
func BuscarIDsFavoritos(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/favoritos/ids", config.APIURL)
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

	var produtoIDs []uint64
	if erro = json.NewDecoder(response.Body).Decode(&produtoIDs); erro != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	respostas.JSON(w, http.StatusOK, produtoIDs)
}

// ToggleFavorito alterna o status de favorito (adiciona ou remove)
func ToggleFavorito(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	produtoIDStr := parametros["produtoId"]

	produtoID, erro := strconv.ParseUint(produtoIDStr, 10, 64)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: "ID do produto inválido"})
		return
	}

	// Primeiro verifica se já é favorito
	urlVerificar := fmt.Sprintf("%s/favoritos/verificar/%s", config.APIURL, produtoIDStr)
	responseVerificar, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, urlVerificar, nil)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer responseVerificar.Body.Close()

	var resultado map[string]bool
	if erro = json.NewDecoder(responseVerificar.Body).Decode(&resultado); erro != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	isFavorito := resultado["isFavorito"]

	// Se já é favorito, remove. Senão, adiciona.
	var response *http.Response
	if isFavorito {
		// Remover
		url := fmt.Sprintf("%s/favoritos/%s", config.APIURL, produtoIDStr)
		response, erro = requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodDelete, url, nil)
	} else {
		// Adicionar
		favoritoJSON, _ := json.Marshal(map[string]uint64{"produtoId": produtoID})
		url := fmt.Sprintf("%s/favoritos", config.APIURL)
		response, erro = requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPost, url, favoritoJSON)
	}

	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	mensagem := "Produto adicionado aos favoritos"
	if isFavorito {
		mensagem = "Produto removido dos favoritos"
	}

	respostas.JSON(w, http.StatusOK, map[string]interface{}{
		"mensagem":   mensagem,
		"isFavorito": !isFavorito,
	})
}
