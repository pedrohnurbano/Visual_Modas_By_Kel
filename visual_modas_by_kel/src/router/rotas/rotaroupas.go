package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotaPaginaRoupas = Rota{
	URI:                "/roupas",
	Metodo:             http.MethodGet,
	Funcao:             controllers.CarregarPaginaRoupas,
	RequerAutenticacao: true,
	RequerAdmin:        false,
}
