package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotasFavoritos = []Rota{
	{
		URI:                "/api/favoritos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.AdicionarFavorito,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	{
		URI:                "/api/favoritos/{produtoId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.RemoverFavorito,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	{
		URI:                "/api/favoritos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarFavoritos,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	{
		URI:                "/api/favoritos/ids",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarIDsFavoritos,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	{
		URI:                "/api/favoritos/toggle/{produtoId}",
		Metodo:             http.MethodPost,
		Funcao:             controllers.ToggleFavorito,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
}
