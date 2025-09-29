package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasProdutos = []Rota{
	// Criar produto
	{
		URI:                "/produtos",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarProduto,
		RequerAutenticacao: true,
	},
	// Buscar todos os produtos (com filtro opcional)
	{
		URI:                "/produtos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProdutos,
		RequerAutenticacao: false, // Permite visualizar sem login
	},
	// Buscar produto específico
	{
		URI:                "/produtos/{produtoId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProduto,
		RequerAutenticacao: false,
	},
	// Buscar produtos por categoria
	{
		URI:                "/produtos/categoria/{categoria}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarProdutosPorCategoria,
		RequerAutenticacao: false,
	},
	// Atualizar produto
	{
		URI:                "/produtos/{produtoId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarProduto,
		RequerAutenticacao: true,
	},
	// Deletar produto (soft delete)
	{
		URI:                "/produtos/{produtoId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarProduto,
		RequerAutenticacao: true,
	},
	// Buscar produtos do usuário logado
	{
		URI:                "/meus-produtos",
		Metodo:             http.MethodGet,
		Funcao:             controllers.MeusProdutos,
		RequerAutenticacao: true,
	},
}