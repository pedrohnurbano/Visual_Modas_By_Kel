package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasUsuarios = []Rota{ //Definindo rotas
	{
		URI:                "/usuarios",              //Endereço
		Metodo:             http.MethodPost,          //POST
		Funcao:             controllers.CriarUsuario, //Função executada
		RequerAutenticacao: false,                    //Não precisa login
	},
	{
		URI:                "/usuarios",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarUsuarios,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/dados",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarDadosUsuarioAutenticado,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{usuarioId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarUsuario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{usuarioId}",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarUsuario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{usuarioId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarUsuario,
		RequerAutenticacao: true,
	},
	{
		URI:                "/usuarios/{usuarioId}/atualizar-senha",
		Metodo:             http.MethodPost,
		Funcao:             controllers.AtualizarSenha,
		RequerAutenticacao: true,
	},
}
