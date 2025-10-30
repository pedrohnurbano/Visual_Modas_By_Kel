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
	// Confirmar entrega do pedido (usuário)
	{
		URI:                "/api/pedidos/{pedidoId}/confirmar-entrega",
		Metodo:             http.MethodPut,
		Funcao:             controllers.ConfirmarEntrega,
		RequerAutenticacao: true,
	},
	// Confirmar pagamento e criar pedido
	{
		URI:                "/api/pedidos/confirmar-pagamento",
		Metodo:             http.MethodPost,
		Funcao:             controllers.ConfirmarPagamentoPedido,
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
	// Atualizar código de rastreio do pedido (admin)
	{
		URI:                "/api/admin/pedidos/{pedidoId}/rastreio",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarCodigoRastreio,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
}
