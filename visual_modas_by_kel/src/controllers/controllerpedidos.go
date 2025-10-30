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

// CriarPedido chama a API para criar um novo pedido
func CriarPedido(w http.ResponseWriter, r *http.Request) {
	body, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: "Erro ao ler dados"})
		return
	}
	defer r.Body.Close()

	url := fmt.Sprintf("%s/pedidos", config.APIURL)
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

// BuscarPedido chama a API para buscar um pedido específico
func BuscarPedido(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	pedidoID := parametros["pedidoId"]

	url := fmt.Sprintf("%s/pedidos/%s", config.APIURL, pedidoID)
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

	var pedido interface{}
	if erro = json.NewDecoder(response.Body).Decode(&pedido); erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: "Erro ao processar resposta"})
		return
	}

	respostas.JSON(w, http.StatusOK, pedido)
}

// BuscarPedidosUsuario chama a API para buscar todos os pedidos do usuário
func BuscarPedidosUsuario(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/pedidos", config.APIURL)
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

	var pedidos interface{}
	if erro = json.NewDecoder(response.Body).Decode(&pedidos); erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: "Erro ao processar resposta"})
		return
	}

	respostas.JSON(w, http.StatusOK, pedidos)
}

// ListarTodosPedidos chama a API para listar todos os pedidos (admin)
func ListarTodosPedidos(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/admin/pedidos", config.APIURL)
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

	var pedidos interface{}
	if erro = json.NewDecoder(response.Body).Decode(&pedidos); erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: "Erro ao processar resposta"})
		return
	}

	respostas.JSON(w, http.StatusOK, pedidos)
}

// AtualizarStatusPedido chama a API para atualizar o status de um pedido (admin)
func AtualizarStatusPedido(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	pedidoID := parametros["pedidoId"]

	body, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: "Erro ao ler dados"})
		return
	}
	defer r.Body.Close()

	url := fmt.Sprintf("%s/admin/pedidos/%s/status", config.APIURL, pedidoID)
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

// AtualizarCodigoRastreio chama a API para atualizar o código de rastreio de um pedido (admin)
func AtualizarCodigoRastreio(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	pedidoID := parametros["pedidoId"]

	body, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: "Erro ao ler dados"})
		return
	}
	defer r.Body.Close()

	url := fmt.Sprintf("%s/admin/pedidos/%s/rastreio", config.APIURL, pedidoID)
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

// ConfirmarEntrega chama a API para confirmar a entrega de um pedido (usuário)
func ConfirmarEntrega(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	pedidoID := parametros["pedidoId"]

	url := fmt.Sprintf("%s/pedidos/%s/confirmar-entrega", config.APIURL, pedidoID)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPut, url, nil)
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

// ConfirmarPagamentoPedido chama a API para confirmar pagamento e criar pedido
func ConfirmarPagamentoPedido(w http.ResponseWriter, r *http.Request) {
	body, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: "Erro ao ler dados"})
		return
	}
	defer r.Body.Close()

	url := fmt.Sprintf("%s/pedidos/confirmar-pagamento", config.APIURL)
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
