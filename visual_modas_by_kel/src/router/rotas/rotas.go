package rotas

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Rota representa todas as rotas da aplicação web
type Rota struct {
	URI                string
	Metodo             string
	Funcao             func(http.ResponseWriter, *http.Request)
	RequerAutenticacao bool
}

// Configurar coloca todas as rotas dentro do router
func Configurar(router *mux.Router) *mux.Router {
	rotas := rotasLogin
	rotas = append(rotas, rotasUsuarios...)
<<<<<<< Updated upstream
=======
	rotas = append(rotas, rotaPaginaPrincipal)
	rotas = append(rotas, rotasAdmin...)
	rotas = append(rotas, rotaLogout)
	rotas = append(rotas, rotasProdutos...)
>>>>>>> Stashed changes

	for _, rota := range rotas {
		router.HandleFunc(rota.URI, rota.Funcao).Methods(rota.Metodo)
	}

	// Serve arquivos estáticos (CSS, JS, imagens)
	fileServer := http.FileServer(http.Dir("./"))
	router.PathPrefix("/").Handler(fileServer)

	return router
}