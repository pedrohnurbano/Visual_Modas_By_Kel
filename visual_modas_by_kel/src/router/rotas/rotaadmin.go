package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/controllers"
)

var rotasAdmin = []Rota{
	{
		URI:                "/painel-admin",
		Metodo:             http.MethodGet,
		Funcao:             controllers.CarregarPainelAdmin,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
	{
		URI:                "/admin/usuarios",
		Metodo:             http.MethodGet,
		Funcao:             controllers.BuscarUsuariosAdmin,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
	{
		URI:                "/admin/usuarios/{usuarioId}/role",
		Metodo:             http.MethodPut,
		Funcao:             controllers.AtualizarRoleUsuario,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
	{
		URI:                "/admin/usuarios/{usuarioId}",
		Metodo:             http.MethodDelete,
		Funcao:             controllers.DeletarUsuarioAdmin,
		RequerAutenticacao: true,
		RequerAdmin:        true,
	},
}
