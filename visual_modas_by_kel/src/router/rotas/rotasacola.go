package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

// Página da sacola - USUÁRIO AUTENTICADO
var rotaPaginaSacola = Rota{
	URI:                "/sacola",
	Metodo:             http.MethodGet,
	Funcao:             controllers.CarregarPaginaSacola,
	RequerAutenticacao: true, // Requer login
	RequerAdmin:        false,
}
