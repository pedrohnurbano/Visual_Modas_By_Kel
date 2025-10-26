package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasCarrinho = []Rota{
	// Adicionar produto ao carrinho
	{
		URI:                "/carrinho",
		Metodo:             http.MethodPost,
		Funcao:             controllers.AdicionarAoCarrinho,
		RequerAutenticacao: true,
	},
	// Buscar todos os itens do carrinho
	{
		URI:                "/carrinho",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarCarrinho,
		RequerAutenticacao: true,
	},
	// Buscar resumo do carrinho com totais
	{
		URI:                "/carrinho/resumo",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarResumoCarrinho,
		RequerAutenticacao: true,
	},
	// Contar itens no carrinho
	{
		URI:                "/carrinho/contar",
		Metodo:             http.MethodGet,
		Funcao:             controllers.ContarItensCarrinho,
		RequerAutenticacao: true,
	},
	// Atualizar quantidade de um item no carrinho
	{
		URI:                "/carrinho/{produtoId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarQuantidadeCarrinho,
		RequerAutenticacao: true,
	},
	// Remover produto do carrinho
	{
		URI:                "/carrinho/{produtoId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.RemoverDoCarrinho,
		RequerAutenticacao: true,
	},
	// Limpar carrinho
	{
		URI:                "/carrinho/limpar",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.LimparCarrinho,
		RequerAutenticacao: true,
	},
}
