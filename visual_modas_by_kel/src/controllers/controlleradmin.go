package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"visual_modas_by_kel/visual_modas_by_kel/src/config"
	"visual_modas_by_kel/visual_modas_by_kel/src/cookies"
	"visual_modas_by_kel/visual_modas_by_kel/src/modelos"
	"visual_modas_by_kel/visual_modas_by_kel/src/requisicoes"
	"visual_modas_by_kel/visual_modas_by_kel/src/respostas"
	"visual_modas_by_kel/visual_modas_by_kel/src/utils"

	"github.com/gorilla/mux"
)

// CarregarPainelAdmin carrega a página do painel administrativo
func CarregarPainelAdmin(w http.ResponseWriter, r *http.Request) {
	cookie, erro := cookies.Ler(r)
	if erro != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Verificar se o usuário é admin
	if cookie["role"] != "admin" {
		http.Redirect(w, r, "/home", http.StatusForbidden)
		return
	}

	// Buscar dados do dashboard
	url := fmt.Sprintf("%s/admin/dashboard", config.APIURL)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if erro != nil {
		utils.ExecutarTemplate(w, "painel-admin.html", nil)
		return
	}
	defer response.Body.Close()

	var dashboard map[string]interface{}
	if response.StatusCode == http.StatusOK {
		if erro = json.NewDecoder(response.Body).Decode(&dashboard); erro != nil {
			dashboard = nil
		}
	}

	utils.ExecutarTemplate(w, "painel-admin.html", dashboard)
}

// BuscarUsuariosAdmin busca todos os usuários (admin only)
func BuscarUsuariosAdmin(w http.ResponseWriter, r *http.Request) {
	cookie, erro := cookies.Ler(r)
	if erro != nil {
		respostas.JSON(w, http.StatusUnauthorized, respostas.ErroAPI{Erro: "Não autorizado"})
		return
	}

	if cookie["role"] != "admin" {
		respostas.JSON(w, http.StatusForbidden, respostas.ErroAPI{Erro: "Acesso negado"})
		return
	}

	url := fmt.Sprintf("%s/admin/usuarios", config.APIURL)
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

	var usuarios []modelos.Usuario
	if erro = json.NewDecoder(response.Body).Decode(&usuarios); erro != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

// AtualizarRoleUsuario atualiza a role de um usuário (admin only)
func AtualizarRoleUsuario(w http.ResponseWriter, r *http.Request) {
	cookie, erro := cookies.Ler(r)
	if erro != nil {
		respostas.JSON(w, http.StatusUnauthorized, respostas.ErroAPI{Erro: "Não autorizado"})
		return
	}

	if cookie["role"] != "admin" {
		respostas.JSON(w, http.StatusForbidden, respostas.ErroAPI{Erro: "Acesso negado"})
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: "ID do usuário inválido"})
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.JSON(w, http.StatusUnprocessableEntity, respostas.ErroAPI{Erro: erro.Error()})
		return
	}

	url := fmt.Sprintf("%s/admin/usuarios/%d/role", config.APIURL, usuarioID)
	response, erro := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPut, url, corpoRequisicao)
	if erro != nil {
		respostas.JSON(w, http.StatusInternalServerError, respostas.ErroAPI{Erro: erro.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		respostas.TratarStatusCodeDeErro(w, response)
		return
	}

	respostas.JSON(w, http.StatusOK, nil)
}

// DeletarUsuarioAdmin deleta um usuário (admin only)
func DeletarUsuarioAdmin(w http.ResponseWriter, r *http.Request) {
	cookie, erro := cookies.Ler(r)
	if erro != nil {
		respostas.JSON(w, http.StatusUnauthorized, respostas.ErroAPI{Erro: "Não autorizado"})
		return
	}

	if cookie["role"] != "admin" {
		respostas.JSON(w, http.StatusForbidden, respostas.ErroAPI{Erro: "Acesso negado"})
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.JSON(w, http.StatusBadRequest, respostas.ErroAPI{Erro: "ID do usuário inválido"})
		return
	}

	url := fmt.Sprintf("%s/admin/usuarios/%d", config.APIURL, usuarioID)
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

	respostas.JSON(w, http.StatusOK, nil)
}
