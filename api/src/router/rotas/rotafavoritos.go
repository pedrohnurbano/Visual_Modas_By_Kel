package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasFavoritos = []Rota{
	// Adicionar produto aos favoritos
	{
		URI:                "/favoritos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.AdicionarFavorito,
		RequerAutenticacao: true,
	},
	// Remover produto dos favoritos
	{
		URI:                "/favoritos/{produtoId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.RemoverFavorito,
		RequerAutenticacao: true,
	},
	// Buscar todos os favoritos do usuário (com dados completos dos produtos)
	{
		URI:                "/favoritos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarFavoritos,
		RequerAutenticacao: true,
	},
	// Buscar apenas os IDs dos produtos favoritos do usuário
	{
		URI:                "/favoritos/ids",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarIDsFavoritos,
		RequerAutenticacao: true,
	},
	// Verificar se um produto está nos favoritos
	{
		URI:                "/favoritos/verificar/{produtoId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.VerificarFavorito,
		RequerAutenticacao: true,
	},
	// Toggle favorito (adicionar ou remover)
	{
		URI:                "/favoritos/toggle/{produtoId}",
		Metodo:             http.MethodPost,
		Funcao:             controllers.ToggleFavorito,
		RequerAutenticacao: true,
	},
}
