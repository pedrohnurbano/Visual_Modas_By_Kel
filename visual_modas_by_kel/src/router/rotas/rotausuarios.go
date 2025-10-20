package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotasUsuarios = []Rota{
	{
		URI:                "/criar-usuario",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPaginaDeCadastroDeUsuario,
		RequerAutenticacao: false,
	},
	{
		URI:                "/usuarios",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarUsuario,
		RequerAutenticacao: false,
	},
	{
		URI:                "/usuario",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPaginaDeUsuario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/api/usuarios/dados",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarDadosUsuario,
		RequerAutenticacao: true,
	},
}
