package rotas

import (
	"net/http"
	"strings"
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
	rotas = append(rotas, rotaPaginaFavoritos)
	rotas = append(rotas, rotaPaginaSacola)
	rotas = append(rotas, rotaPaginaRoupas)
	rotas = append(rotas, rotaPaginaCheckout)
	rotas = append(rotas, rotaPaginaCheckoutConfirmacao)
	rotas = append(rotas, rotaPaginaTermosPolitica)
	rotas = append(rotas, rotasAdmin...)
	rotas = append(rotas, rotaLogout)
	rotas = append(rotas, rotasProdutos...)
	rotas = append(rotas, rotasFavoritos...)
	rotas = append(rotas, rotasCarrinho...)
	rotas = append(rotas, rotasPedidos...)
	rotas = append(rotas, rotaAbacatePay)

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

	// Serve arquivos estáticos (CSS, JS, imagens) - APENAS para arquivos que não são rotas da API
	fileServer := http.FileServer(http.Dir("./"))
	router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar se é uma rota da API
		if strings.HasPrefix(r.URL.Path, "/api/") {
			// Se for rota da API, não servir arquivo estático
			http.NotFound(w, r)
			return
		}
		// Não servir arquivos .html diretamente - devem usar as rotas definidas
		if strings.HasSuffix(r.URL.Path, ".html") {
			http.NotFound(w, r)
			return
		}
		// Para outras rotas, servir arquivo estático
		fileServer.ServeHTTP(w, r)
	}))

	return router
}
