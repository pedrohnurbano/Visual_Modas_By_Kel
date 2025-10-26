package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotasPedidos = []Rota{
	// Criar novo pedido
	{
		URI:                "/api/pedidos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarPedido,
		RequerAutenticacao: true,
	},
	// Buscar pedidos do usuário logado
	{
		URI:                "/api/pedidos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPedidosUsuario,
		RequerAutenticacao: true,
	},
	// Buscar pedido específico pelo ID
	{
		URI:                "/api/pedidos/{pedidoId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarPedido,
		RequerAutenticacao: true,
	},
	// Listar todos os pedidos (admin)
	{
		URI:                "/api/admin/pedidos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.ListarTodosPedidos,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
	// Atualizar status do pedido (admin)
	{
		URI:                "/api/admin/pedidos/{pedidoId}/status",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarStatusPedido,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
}
