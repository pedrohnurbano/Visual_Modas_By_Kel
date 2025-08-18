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
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// CriarUsuario insere um usuário no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := io.ReadAll(r.Body) //Lê o corpo da requisição
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro) //recebe 3 parâmetros: resposta, código de status e mensagem de erro
		return
	}

	var usuario modelos.Usuario //criando usuario que está no pacote de modelos
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro) //recebe 3 parâmetros: resposta, código de status e mensagem de erro
		return
	}

	if erro = usuario.Preparar("cadastro"); erro != nil { //recebe o cadastro como parametro, para quando chegar no método de validar, ver que a etapa == cadastro, e se a senha estiver em branco vai fazer o erro q falta senha
		respostas.Erro(w, http.StatusBadRequest, erro) //recebe 3 parâmetros: resposta, código de status e mensagem de erro
		return
	}

	db, erro := banco.Conectar() //criando conexao com o banco de dados
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro) //recebe 3 parâmetros: resposta, código de status e mensagem de erro
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db) //criando um novo repositorio de usuarios

	// Verificar se email já existe
	_, erro = repositorio.BuscarPorEmail(usuario.Email)
	if erro == nil {
		respostas.Erro(w, http.StatusConflict, errors.New("este e-mail já está cadastrado"))
		return
	}

	// Verificar se CPF já existe
	_, erro = repositorio.BuscarPorCPF(usuario.CPF)
	if erro == nil {
		respostas.Erro(w, http.StatusConflict, errors.New("este CPF já está cadastrado"))
		return
	}

	usuario.ID, erro = repositorio.Criar(usuario) //chamando o metodo criar do repositorio de usuarios (inserindo)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro) //recebe 3 parâmetros: resposta, código de status e mensagem de erro
		return
	}

	// Limpar senha da resposta por segurança
	usuario.Senha = ""
	respostas.JSON(w, http.StatusCreated, usuario) //enviando resposta para o usuario
}

// BuscarUsuarios busca todos os usuários salvos no banco (apenas para administradores)
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuEmail := strings.ToLower(r.URL.Query().Get("usuario")) //converte a string para minuscula, vai trazer tudo que tiver na query (rota) e vai pegar o valor do campo usuario
	db, erro := banco.Conectar()                                 //conecta ao banco de dados
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}
	defer db.Close() //fecha a conexão com o banco de dados

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.Buscar(nomeOuEmail)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	// Limpar senhas da resposta por segurança
	for i := range usuarios {
		usuarios[i].Senha = ""
	}

	respostas.JSON(w, http.StatusOK, usuarios)
}

// BuscarUsuario busca um usuário salvo no banco
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r) //mux.Vars recebe a requisição e retorna um map com os parâmetros da rota

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64) //converte o id do usuario de string para uint64
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Só permite que o usuário veja seus próprios dados
	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível visualizar dados de outro usuário"))
		return
	}

	db, erro := banco.Conectar() //conecta ao banco de dados
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db) //criando um novo repositorio de usuarios
	usuario, erro := repositorio.BuscarPorID(usuarioID)       //chamando o metodo buscar por id do repositorio de usuarios, passando o id do usuario como parâmetro
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	// Limpar senha da resposta por segurança
	usuario.Senha = ""
	respostas.JSON(w, http.StatusOK, usuario)
}

// AtualizarUsuario altera as informações de um usuário no banco
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r) //mux.Vars recebe a requisição e retorna um map com os parâmetros da rota

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64) //converte o id do usuario de string para uint64
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	if usuarioID != usuarioIDNoToken { //nao deixa os usuarios alterarem informações que não são deles (dos outros)
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar um usuário que não seja o seu")) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body) //lê o corpo da requisição
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil { //json.Unmarshall converte o corpo da requisição em um objeto do tipo usuario
		respostas.Erro(w, http.StatusBadRequest, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	// Verificar se o email já existe para outro usuário
	usuarioComEmail, erro := repositorio.BuscarPorEmail(usuario.Email)
	if erro == nil && usuarioComEmail.ID != usuarioID {
		respostas.Erro(w, http.StatusConflict, errors.New("este e-mail já está sendo usado por outro usuário"))
		return
	}

	// Verificar se o CPF já existe para outro usuário
	usuarioComCPF, erro := repositorio.BuscarPorCPF(usuario.CPF)
	if erro == nil && usuarioComCPF.ID != usuarioID {
		respostas.Erro(w, http.StatusConflict, errors.New("este CPF já está sendo usado por outro usuário"))
		return
	}

	if erro = repositorio.Atualizar(usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro) //chama a função erro lá do respostas.go e envia o erro para o usuario
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}

// DeletarUsuario exclui as informações de um usuário do banco
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDNoToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível deletar um usuário que não seja o seu"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro) //chama a função erro
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

func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if usuarioIDNoToken != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("nao é possivel atualizar a senha de um usuario que nao seja o seu"))
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	var senha modelos.Senha
	if erro = json.Unmarshal(corpoRequisicao, &senha); erro != nil {
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
	senhaSalvaNoBanco, erro := repositorio.BuscarSenha(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(senhaSalvaNoBanco, senha.Atual); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("a senha atual nao condiz com a que está salva no banco"))
		return
	}

	senhaComHash, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.AtualizarSenha(usuarioID, string(senhaComHash)); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusNoContent, nil)
}
