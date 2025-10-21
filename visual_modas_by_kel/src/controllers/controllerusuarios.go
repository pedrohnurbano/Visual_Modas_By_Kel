package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/config"
	"visual_modas_by_kel/visual_modas_by_kel/src/requisicoes"
	"visual_modas_by_kel/visual_modas_by_kel/src/respostas"
	"visual_modas_by_kel/visual_modas_by_kel/src/utils"

	"github.com/gorilla/mux"
)

// CriarUsuario chama a API para cadastrar um usuário no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	usuario, erro := json.Marshal(map[string]string{
		"nome":      r.FormValue("nome"),
		"sobrenome": r.FormValue("sobrenome"),
		"email":     r.FormValue("email"),
		"senha":     r.FormValue("senha"),
		"telefone":  r.FormValue("telefone"),
		"cpf":       r.FormValue("cpf"),
	})

	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/usuarios", config.APIURL)
	response, erro := http.Post(url, "/application/json", bytes.NewBuffer(usuario))

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

// CarregarPaginaDeUsuario carrega a página de perfil do usuário
func CarregarPaginaDeUsuario(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "usuario.html", nil)
}

// BuscarDadosUsuario chama a API para buscar os dados do usuário autenticado
func BuscarDadosUsuario(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/usuarios/dados", config.APIURL)
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

	var usuario map[string]interface{}
	if erro := json.NewDecoder(response.Body).Decode(&usuario); erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	respostas.JSON(w, http.StatusOK, usuario)
}

// AtualizarDadosUsuario chama a API para atualizar os dados do usuário
func AtualizarDadosUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID := parametros["usuarioId"]

	var dadosUsuario map[string]interface{}
	if erro := json.NewDecoder(r.Body).Decode(&dadosUsuario); erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	dadosJSON, erro := json.Marshal(dadosUsuario)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/usuarios/%s", config.APIURL, usuarioID)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPut, url, dadosJSON)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	// Não tenta ler o body se for 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// AtualizarSenhaUsuario chama a API para atualizar a senha do usuário
func AtualizarSenhaUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID := parametros["usuarioId"]

	var senhaData map[string]interface{}
	if erro := json.NewDecoder(r.Body).Decode(&senhaData); erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	senhaJSON, erro := json.Marshal(senhaData)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/usuarios/%s/atualizar-senha", config.APIURL, usuarioID)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPost, url, senhaJSON)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	// Não tenta ler o body se for 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// DeletarUsuario chama a API para excluir a conta do usuário
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID := parametros["usuarioId"]

	var dados map[string]interface{}
	if erro := json.NewDecoder(r.Body).Decode(&dados); erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	dadosJSON, erro := json.Marshal(dados)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/usuarios/%s", config.APIURL, usuarioID)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodDelete, url, dadosJSON)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	// Não tenta ler o body se for 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
