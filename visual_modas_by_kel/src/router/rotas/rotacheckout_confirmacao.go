package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

// Rota para a página de confirmação do checkout (passo 5)
var rotaPaginaCheckoutConfirmacao = Rota{
	URI:                "/checkout/confirmacao",
	Metodo:             http.MethodGet,
	Funcao:             controllers.CarregarPaginaCheckout,
	RequerAutenticacao: true,
	RequerAdmin:        false,
}
