package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotaAbacatePay = Rota{
	URI:                "/api/abacatepay/criar-cobranca",
	Metodo:             http.MethodPost,
	Funcao:             controllers.CriarCobrancaAbacatePay,
	RequerAutenticacao: true,
	RequerAdmin:        false,
}
