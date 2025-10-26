package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotasCarrinho = []Rota{
	// Adicionar produto ao carrinho
	{
		URI:                "/api/carrinho",
		Metodo:             http.MethodPost,
		Funcao:             controllers.AdicionarAoCarrinho,
		RequerAutenticacao: true,
	},
	// Buscar todos os itens do carrinho
	{
		URI:                "/api/carrinho",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarCarrinho,
		RequerAutenticacao: true,
	},
	// Buscar resumo do carrinho com totais
	{
		URI:                "/api/carrinho/resumo",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarResumoCarrinho,
		RequerAutenticacao: true,
	},
	// Contar itens no carrinho
	{
		URI:                "/api/carrinho/contar",
		Metodo:             http.MethodGet,
		Funcao:             controllers.ContarItensCarrinho,
		RequerAutenticacao: true,
	},
	// Atualizar quantidade de um item no carrinho
	{
		URI:                "/api/carrinho/{produtoId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarQuantidadeCarrinho,
		RequerAutenticacao: true,
	},
	// Remover produto do carrinho
	{
		URI:                "/api/carrinho/{produtoId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.RemoverDoCarrinho,
		RequerAutenticacao: true,
	},
	// Limpar carrinho
	{
		URI:                "/api/carrinho/limpar",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.LimparCarrinho,
		RequerAutenticacao: true,
	},
}
