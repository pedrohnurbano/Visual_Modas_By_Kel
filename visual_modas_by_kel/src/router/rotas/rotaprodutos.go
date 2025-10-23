package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotasProdutos = []Rota{
	{
		URI:                "/cadastro-produto",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPaginaCadastroProduto,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
	{
		URI:                "/produtos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPaginaProdutos,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	{
		URI:                "/meus-produtos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPaginaMeusProdutos,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
	{
		URI:                "/api/produtos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarProduto,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
	{
		URI:                "/api/produtos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProdutos,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	{
		URI:                "/api/produtos/home",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProdutosHome,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	{
		URI:                "/api/produtos/secao/{secao}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProdutosPorSecao,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	{
		URI:                "/api/produtos/{produtoId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProduto,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	{
		URI:                "/api/meus-produtos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarMeusProdutos,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
	{
		URI:                "/api/produtos/{produtoId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarProduto,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
	{
		URI:                "/api/produtos/{produtoId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarProduto,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
}
