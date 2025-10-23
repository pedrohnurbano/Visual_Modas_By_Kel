package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotaPaginaCheckout = Rota{
	URI:                "/checkout",
	Metodo:             http.MethodGet,
	Funcao:             controllers.CarregarPaginaCheckout,
	RequerAutenticacao: true,
	RequerAdmin:        false,
}
