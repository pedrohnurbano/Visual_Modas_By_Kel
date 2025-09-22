package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// verificarAdmin verifica se o usuário tem role de admin
func verificarAdmin(r *http.Request) error {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		return erro
	}

	db, erro := banco.Conectar()
	if erro != nil {
		return erro
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorio.BuscarPorID(usuarioID)
	if erro != nil {
		return erro
	}

	if usuario.Role != "admin" {
		return errors.New("acesso negado: apenas administradores podem acessar esta funcionalidade")
	}

	return nil
}

// ListarTodosUsuarios lista todos os usuários (apenas admin)
func ListarTodosUsuarios(w http.ResponseWriter, r *http.Request) {
	// Verifica se é admin
	if erro := verificarAdmin(r); erro != nil {
		respostas.Erro(w, http.StatusForbidden, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.Buscar("") // Busca todos
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Limpar senhas da resposta por segurança
	for i := range usuarios {
		usuarios[i].Senha = ""
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

// BuscarUsuarioAdmin busca qualquer usuário (apenas admin)
func BuscarUsuarioAdmin(w http.ResponseWriter, r *http.Request) {
	// Verifica se é admin
	if erro := verificarAdmin(r); erro != nil {
		respostas.Erro(w, http.StatusForbidden, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorio.BuscarPorID(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Limpar senha da resposta por segurança
	usuario.Senha = ""
	respostas.JSON(w, http.StatusOK, usuario)
}

// AtualizarRoleUsuario atualiza a role de um usuário (apenas admin)
func AtualizarRoleUsuario(w http.ResponseWriter, r *http.Request) {
	// Verifica se é admin
	if erro := verificarAdmin(r); erro != nil {
		respostas.Erro(w, http.StatusForbidden, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var dados struct {
		Role string `json:"role"`
	}

	if erro = json.Unmarshal(corpoRequisicao, &dados); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Valida a role
	if dados.Role != "user" && dados.Role != "admin" {
		respostas.Erro(w, http.StatusBadRequest, errors.New("role inválida. Use 'user' ou 'admin'"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.AtualizarRole(usuarioID, dados.Role); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarUsuarioAdmin deleta qualquer usuário (apenas admin)
func DeletarUsuarioAdmin(w http.ResponseWriter, r *http.Request) {
	// Verifica se é admin
	if erro := verificarAdmin(r); erro != nil {
		respostas.Erro(w, http.StatusForbidden, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Previne admin de deletar a si mesmo
	adminID, _ := autenticacao.ExtrairUsuarioID(r)
	if usuarioID == adminID {
		respostas.Erro(w, http.StatusBadRequest, errors.New("você não pode deletar sua própria conta"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Deletar(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DashboardAdmin retorna estatísticas para o admin
func DashboardAdmin(w http.ResponseWriter, r *http.Request) {
	// Verifica se é admin
	if erro := verificarAdmin(r); erro != nil {
		respostas.Erro(w, http.StatusForbidden, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Busca estatísticas básicas
	var stats struct {
		TotalUsuarios int `json:"total_usuarios"`
		TotalAdmins   int `json:"total_admins"`
		TotalUsers    int `json:"total_users"`
	}

	// Total de usuários
	erro = db.QueryRow("SELECT COUNT(*) FROM usuarios").Scan(&stats.TotalUsuarios)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Total de admins
	erro = db.QueryRow("SELECT COUNT(*) FROM usuarios WHERE role = 'admin'").Scan(&stats.TotalAdmins)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Total de users comuns
	erro = db.QueryRow("SELECT COUNT(*) FROM usuarios WHERE role = 'user'").Scan(&stats.TotalUsers)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, stats)
}
