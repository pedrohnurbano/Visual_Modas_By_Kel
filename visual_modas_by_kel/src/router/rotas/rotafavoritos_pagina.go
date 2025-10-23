package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotaPaginaFavoritos = Rota{
	URI:                "/favoritos",
	Metodo:             http.MethodGet,
	Funcao:             controllers.CarregarPaginaFavoritos,
	RequerAutenticacao: true,
	RequerAdmin:        false,
}
