package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/modelos"
	"api/src/repositorios"
	"api/src/respostas"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CriarPedido cria um novo pedido a partir do carrinho do usuário
func CriarPedido(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do usuário do token
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Ler corpo da requisição
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	var dados modelos.CriarPedidoRequest
	if erro = json.Unmarshal(corpoRequisicao, &dados); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Validar dados obrigatórios
	if dados.NomeCompleto == "" || dados.Email == "" || dados.Telefone == "" ||
		dados.Endereco == "" || dados.Numero == "" || dados.Bairro == "" ||
		dados.Cidade == "" || dados.Estado == "" || dados.CEP == "" {
		respostas.Erro(w, http.StatusBadRequest, fmt.Errorf("todos os campos obrigatórios devem ser preenchidos"))
		return
	}

	// Conectar ao banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Buscar resumo do carrinho para calcular o total
	repositorioCarrinho := repositorios.NovoRepositorioDeCarrinho(db)
	resumo, erro := repositorioCarrinho.BuscarResumoCarrinho(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Verificar se o carrinho está vazio
	if len(resumo.Itens) == 0 {
		respostas.Erro(w, http.StatusBadRequest, fmt.Errorf("carrinho vazio"))
		return
	}

	// Criar pedido
	pedido := modelos.Pedido{
		UsuarioID:      usuarioID,
		NomeCompleto:   dados.NomeCompleto,
		Email:          dados.Email,
		Telefone:       dados.Telefone,
		Endereco:       dados.Endereco,
		Numero:         dados.Numero,
		Complemento:    dados.Complemento,
		Bairro:         dados.Bairro,
		Cidade:         dados.Cidade,
		Estado:         dados.Estado,
		CEP:            dados.CEP,
		FormaPagamento: dados.FormaPagamento,
		Status:         "pendente",
		Total:          resumo.ValorTotal,
	}

	repositorioPedidos := repositorios.NovoRepositorioDePedidos(db)
	pedidoID, erro := repositorioPedidos.Criar(pedido)
	if erro != nil {
		println("ERRO ao criar pedido:", erro.Error())
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	println("Pedido criado com sucesso! PedidoID:", pedidoID, "UsuarioID:", usuarioID)
	respostas.JSON(w, http.StatusCreated, map[string]interface{}{
		"mensagem": "Pedido criado com sucesso",
		"pedidoId": pedidoID,
	})
}

// BuscarPedido busca um pedido específico pelo ID
func BuscarPedido(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do usuário do token
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Extrair ID do pedido dos parâmetros da URL
	parametros := mux.Vars(r)
	pedidoID, erro := strconv.ParseUint(parametros["pedidoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Conectar ao banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Buscar pedido completo
	repositorio := repositorios.NovoRepositorioDePedidos(db)
	pedidoCompleto, erro := repositorio.BuscarPedidoCompleto(pedidoID)
	if erro != nil {
		respostas.Erro(w, http.StatusNotFound, erro)
		return
	}

	// Verificar se o pedido pertence ao usuário
	if pedidoCompleto.Pedido.UsuarioID != usuarioID {
		respostas.Erro(w, http.StatusForbidden, fmt.Errorf("acesso negado"))
		return
	}

	respostas.JSON(w, http.StatusOK, pedidoCompleto)
}

// BuscarPedidosUsuario busca todos os pedidos do usuário logado
func BuscarPedidosUsuario(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do usuário do token
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Conectar ao banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Buscar pedidos do usuário
	repositorio := repositorios.NovoRepositorioDePedidos(db)
	pedidos, erro := repositorio.BuscarPorUsuario(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, pedidos)
}

// ListarTodosPedidos lista todos os pedidos (admin)
func ListarTodosPedidos(w http.ResponseWriter, r *http.Request) {
	// Conectar ao banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Listar todos os pedidos
	repositorio := repositorios.NovoRepositorioDePedidos(db)
	pedidos, erro := repositorio.ListarTodos()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, pedidos)
}

// AtualizarStatusPedido atualiza o status de um pedido (admin)
func AtualizarStatusPedido(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do pedido dos parâmetros da URL
	parametros := mux.Vars(r)
	pedidoID, erro := strconv.ParseUint(parametros["pedidoId"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Ler corpo da requisição
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	var dados struct {
		Status string `json:"status"`
	}
	if erro = json.Unmarshal(corpoRequisicao, &dados); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Validar status
	if dados.Status == "" {
		respostas.Erro(w, http.StatusBadRequest, fmt.Errorf("status não pode ser vazio"))
		return
	}

	// Conectar ao banco
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Atualizar status do pedido
	repositorio := repositorios.NovoRepositorioDePedidos(db)
	if erro = repositorio.AtualizarStatus(pedidoID, dados.Status); erro != nil {
		println("ERRO ao atualizar status do pedido:", erro.Error())
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	println("Status do pedido atualizado! PedidoID:", pedidoID, "Status:", dados.Status)
	respostas.JSON(w, http.StatusOK, map[string]string{"mensagem": "Status do pedido atualizado"})
}
