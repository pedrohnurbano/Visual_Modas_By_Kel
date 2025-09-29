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
	RequerAdmin        bool
}

// Configurar coloca todas as rotas dentro do router
func Configurar(router *mux.Router) *mux.Router {
	rotas := rotasLogin
	rotas = append(rotas, rotasUsuarios...)
	rotas = append(rotas, rotaPaginaPrincipal)
	rotas = append(rotas, rotasAdmin...)
	rotas = append(rotas, rotaLogout)
	rotas = append(rotas, rotasProdutos...) // ADICIONAR ESTA LINHA

	for _, rota := range rotas {
		if rota.RequerAdmin {
			// Rotas que requerem admin
			router.HandleFunc(rota.URI,
				middlewares.Logger(middlewares.AutenticarAdmin(rota.Funcao)),
			).Methods(rota.Metodo)
		} else if rota.RequerAutenticacao {
			// Rotas que requerem apenas autenticação
			router.HandleFunc(rota.URI,
				middlewares.Logger(middlewares.Autenticar(rota.Funcao)),
			).Methods(rota.Metodo)
		} else {
			// Rotas públicas
			router.HandleFunc(rota.URI,
				middlewares.Logger(rota.Funcao),
			).Methods(rota.Metodo)
		}
	}

	// Serve arquivos estáticos (CSS, JS, imagens)
	fileServer := http.FileServer(http.Dir("./"))
	router.PathPrefix("/").Handler(fileServer)

	return router
}