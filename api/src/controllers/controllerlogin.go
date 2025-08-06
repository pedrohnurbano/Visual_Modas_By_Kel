package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"io"
	"net/http"
)

// Login é responsável por autenticar o usuário na API
func Login(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := io.ReadAll(r.Body) //lê o corpo da requisição
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro) //chama a função erro
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioSalvoNoBanco, erro := repositorio.BuscarPorEmail(usuario.Email)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(usuarioSalvoNoBanco.Senha, usuario.Senha); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	token, erro := autenticacao.CriarToken(usuarioSalvoNoBanco.ID) //token é criado pelo id do usuario salvo no banco
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	
	w.Write([]byte(token))
}
