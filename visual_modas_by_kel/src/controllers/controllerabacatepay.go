package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/config"
	"visual_modas_by_kel/visual_modas_by_kel/src/requisicoes"
	"visual_modas_by_kel/visual_modas_by_kel/src/respostas"
)

// CriarCobrancaAbacatePay faz proxy da requisição para a API
func CriarCobrancaAbacatePay(w http.ResponseWriter, r *http.Request) {
	// Ler corpo da requisição
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	// Fazer requisição para a API backend
	url := fmt.Sprintf("%s/abacatepay/criar-cobranca", config.APIURL)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPost, url, corpoRequisicao)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	// Verificar se houve erro na API
	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	// Ler resposta da API
	var resultado map[string]interface{}
	if erro = json.NewDecoder(response.Body).Decode(&resultado); erro != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	// Retornar resposta para o frontend
	respostas.JSON(w, response.StatusCode, resultado)
}

// AbacatePayWebhook faz proxy do webhook para a API
func AbacatePayWebhook(w http.ResponseWriter, r *http.Request) {
	// Ler corpo da requisição
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	// Fazer requisição para a API backend (sem autenticação, pois é webhook externo)
	url := fmt.Sprintf("%s/abacatepay/webhook", config.APIURL)
	req, erro := http.NewRequest(http.MethodPost, url, io.NopCloser(bytes.NewBuffer(corpoRequisicao)))
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, erro := client.Do(req)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	// Retornar resposta
	var resultado map[string]interface{}
	json.NewDecoder(response.Body).Decode(&resultado)
	respostas.JSON(w, response.StatusCode, resultado)
}
