package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotasUsuarios = []Rota{
	{
		URI:                "/cadastro",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPaginaDeCadastroDeUsuario,
		RequerAutenticacao: false,
	},
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
	{
		URI:                "/api/usuarios/{usuarioId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarDadosUsuario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/api/usuarios/{usuarioId}/atualizar-senha",
		Metodo:             http.MethodPost,
		Funcao:             controllers.AtualizarSenhaUsuario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/api/usuarios/{usuarioId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarUsuario,
		RequerAutenticacao: true,
	},
}
