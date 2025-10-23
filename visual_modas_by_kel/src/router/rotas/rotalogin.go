package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotasLogin = []Rota{
	{
		URI:                "/",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarTelaDeLogin,
		RequerAutenticacao: false,
		RequerAdmin:        false,
	},
	{
		URI:                "/login",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarTelaDeLogin,
		RequerAutenticacao: false,
		RequerAdmin:        false,
	},
	{
		URI:                "/login",
		Metodo:             http.MethodPost,
		Funcao:             controllers.FazerLogin,
		RequerAutenticacao: false,
		RequerAdmin:        false,
	},
}
