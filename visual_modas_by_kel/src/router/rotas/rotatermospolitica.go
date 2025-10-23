package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotaPaginaTermosPolitica = Rota{
	URI:                "/termos-politica",
	Metodo:             http.MethodGet,
	Funcao:             controllers.CarregarPaginaTermosPolitica,
	RequerAutenticacao: false,
	RequerAdmin:        false,
}
