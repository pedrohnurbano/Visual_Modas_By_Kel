package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotasAbacatePay = []Rota{
	{
		URI:                "/api/abacatepay/criar-cobranca",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarCobrancaAbacatePay,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	{
		URI:                "/api/abacatepay/webhook",
		Metodo:             http.MethodPost,
		Funcao:             controllers.AbacatePayWebhook,
		RequerAutenticacao: false,
		RequerAdmin:        false,
	},
}
