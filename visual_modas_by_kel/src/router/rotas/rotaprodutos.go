package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotasProdutos = []Rota{
	// Página de cadastro de produto
	{
		URI:                "/cadastro-produto",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPaginaCadastroProduto,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	// Página de listagem de produtos
	{
		URI:                "/produtos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPaginaProdutos,
		RequerAutenticacao: false,
		RequerAdmin:        false,
	},
	// Página meus produtos
	{
		URI:                "/meus-produtos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPaginaMeusProdutos,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	// API: Criar produto
	{
		URI:                "/api/produtos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarProduto,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	// API: Buscar todos produtos
	{
		URI:                "/api/produtos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProdutos,
		RequerAutenticacao: false,
		RequerAdmin:        false,
	},
	// API: Buscar produto específico
	{
		URI:                "/api/produtos/{produtoId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProduto,
		RequerAutenticacao: false,
		RequerAdmin:        false,
	},
	// API: Buscar produtos por seção
	{
		URI:                "/api/produtos/secao/{secao}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProdutosPorSecao,
		RequerAutenticacao: false,
		RequerAdmin:        false,
	},
	// API: Buscar meus produtos
	{
		URI:                "/api/meus-produtos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarMeusProdutos,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	// API: Atualizar produto
	{
		URI:                "/api/produtos/{produtoId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarProduto,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
	// API: Deletar produto
	{
		URI:                "/api/produtos/{produtoId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarProduto,
		RequerAutenticacao: true,
		RequerAdmin:        false,
	},
}