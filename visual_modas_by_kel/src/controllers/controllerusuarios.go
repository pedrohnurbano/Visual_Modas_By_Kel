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
