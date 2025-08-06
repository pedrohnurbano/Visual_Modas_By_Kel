package rotas

import (
	"api/src/controllers"
	"net/http"
)

var rotaLogin = Rota{
	URI:    "/login",
	Metodo: http.MethodPost,
	Funcao: controllers.Login, //Função que será chamada quando a rota for acessada
	RequerAutenticacao: false,
}
