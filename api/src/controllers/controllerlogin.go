package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"api/src/seguranca"
	"encoding/json"
	"errors"
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

	var dadosLogin struct {
		Email string `json:"email"`
		Senha string `json:"senha"`
	}

	if erro = json.Unmarshal(corpoRequisicao, &dadosLogin); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	// Validações básicas
	if dadosLogin.Email == "" {
		respostas.Erro(w, http.StatusBadRequest, errors.New("o e-mail é obrigatório"))
		return
	}

	if dadosLogin.Senha == "" {
		respostas.Erro(w, http.StatusBadRequest, errors.New("a senha é obrigatória"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro) //chama a função erro
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioSalvoNoBanco, erro := repositorio.BuscarPorEmail(dadosLogin.Email)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("credenciais inválidas"))
		return
	}

	// Verificar se o usuário foi encontrado
	if usuarioSalvoNoBanco.ID == 0 {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("credenciais inválidas"))
		return
	}

	if erro = seguranca.VerificarSenha(usuarioSalvoNoBanco.Senha, dadosLogin.Senha); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("credenciais inválidas"))
		return
	}

	token, erro := autenticacao.CriarToken(usuarioSalvoNoBanco.ID) //token é criado pelo id do usuario salvo no banco
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Buscar dados completos do usuário para retornar (sem senha)
	usuarioCompleto, erro := repositorio.BuscarPorID(usuarioSalvoNoBanco.ID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Limpar senha por segurança
	usuarioCompleto.Senha = ""

	// Retornar token e dados do usuário
	resposta := struct {
		Token   string          `json:"token"`
		Usuario modelos.Usuario `json:"usuario"`
	}{
		Token:   token,
		Usuario: usuarioCompleto,
	}

	respostas.JSON(w, http.StatusOK, resposta)
}
