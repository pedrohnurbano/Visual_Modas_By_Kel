package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasPedidos = []Rota{
	// Criar novo pedido
	{
		URI:                "/pedidos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarPedido,
		RequerAutenticacao: true,
	},
	// Buscar pedidos do usuário logado
	{
		URI:                "/pedidos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPedidosUsuario,
		RequerAutenticacao: true,
	},
	// Buscar pedido específico pelo ID
	{
		URI:                "/pedidos/{pedidoId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPedido,
		RequerAutenticacao: true,
	},
	// Listar todos os pedidos (admin)
	{
		URI:                "/admin/pedidos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.ListarTodosPedidos,
		RequerAutenticacao: true,
	},
	// Atualizar status do pedido (admin)
	{
		URI:                "/admin/pedidos/{pedidoId}/status",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarStatusPedido,
		RequerAutenticacao: true,
	},
}
