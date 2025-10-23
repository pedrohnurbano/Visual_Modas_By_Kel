package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotaLogout = Rota{
	URI:                "/logout",
	Metodo:             http.MethodPost,
	Funcao:             controllers.FazerLogout,
	RequerAutenticacao: false,
	RequerAdmin:        false,
}
