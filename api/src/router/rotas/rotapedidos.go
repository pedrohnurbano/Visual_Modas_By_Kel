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
	// Confirmar entrega do pedido (usuário)
	{
		URI:                "/pedidos/{pedidoId}/confirmar-entrega",
		Metodo:             http.MethodPut,
		Funcao:             controllers.ConfirmarEntrega,
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
	// Atualizar código de rastreio do pedido (admin)
	{
		URI:                "/admin/pedidos/{pedidoId}/rastreio",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarCodigoRastreio,
		RequerAutenticacao: true,
	},
	// Criar cobrança no AbacatePay (rota da API)
	{
		URI:                "/abacatepay/criar-cobranca",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarCobrancaAbacatePay,
		RequerAutenticacao: true,
	},
	// Webhook do AbacatePay (não requer autenticação)
	{
		URI:                "/abacatepay/webhook",
		Metodo:             http.MethodPost,
		Funcao:             controllers.AbacatePayWebhook,
		RequerAutenticacao: false,
	},
	// Confirmar pagamento após retorno do AbacatePay
	{
		URI:                "/pedidos/confirmar-pagamento",
		Metodo:             http.MethodPost,
		Funcao:             controllers.ConfirmarPagamentoPedido,
		RequerAutenticacao: true,
	},
}
