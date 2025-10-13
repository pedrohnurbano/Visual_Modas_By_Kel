package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotasFavoritos = []Rota{
	// API: Adicionar produto aos favoritos
	{
		URI:                "/api/favoritos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.AdicionarFavorito,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	// API: Remover produto dos favoritos
	{
		URI:                "/api/favoritos/{produtoId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.RemoverFavorito,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	// API: Buscar todos os favoritos do usuário (com dados completos)
	{
		URI:                "/api/favoritos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarFavoritos,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	// API: Buscar apenas IDs dos favoritos
	{
		URI:                "/api/favoritos/ids",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarIDsFavoritos,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	// API: Toggle favorito (adiciona ou remove)
	{
		URI:                "/api/favoritos/toggle/{produtoId}",
		Metodo:             http.MethodPost,
		Funcao:             controllers.ToggleFavorito,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
}
