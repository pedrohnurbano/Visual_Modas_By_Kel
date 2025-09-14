package rotas

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/middlewares"

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
	rotas = append(rotas, rotaPaginaPrincipal)

	for _, rota := range rotas {

		if rota.RequerAutenticacao {
			router.HandleFunc(rota.URI,
				middlewares.Logger(middlewares.Autenticar(rota.Funcao)),
			).Methods(rota.Metodo)

		} else {
			router.HandleFunc(rota.URI,
				middlewares.Logger(rota.Funcao),
			).Methods(rota.Metodo)
		}
		
		router.HandleFunc(rota.URI, rota.Funcao).Methods(rota.Metodo)
	}

	// Serve arquivos estáticos diretamente da pasta raiz
	fileServer := http.FileServer(http.Dir("./"))
	router.PathPrefix("/").Handler(fileServer)

	return router
}
