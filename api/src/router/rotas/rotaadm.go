package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotasAdmin = []Rota{
	// Rota para listar todos os usuários (admin)
	{
		URI:                "/admin/usuarios",
		Metodo:             http.MethodGet,
		Funcao:             controllers.ListarTodosUsuarios,
		RequerAutenticacao: true,
	},
	// Rota para buscar um usuário específico (admin)
	{
		URI:                "/admin/usuarios/{usuarioId}",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarUsuarioAdmin,
		RequerAutenticacao: true,
	},
	// Rota para atualizar role de um usuário (admin)
	{
		URI:                "/admin/usuarios/{usuarioId}/role",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarRoleUsuario,
		RequerAutenticacao: true,
	},
	// Rota para deletar qualquer usuário (admin)
	{
		URI:                "/admin/usuarios/{usuarioId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarUsuarioAdmin,
		RequerAutenticacao: true,
	},
	// Rota para dashboard do admin
	{
		URI:                "/admin/dashboard",
		Metodo:             http.MethodGet,
		Funcao:             controllers.DashboardAdmin,
		RequerAutenticacao: true,
	},
}
