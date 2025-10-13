package rotas

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// Rota representa todas as rotas da API
type Rota struct {
	URI                string                                   //Endereço da rota
	Metodo             string                                   //Método HTTP
	Funcao             func(http.ResponseWriter, *http.Request) //Função que irá lidar com a requisição
	RequerAutenticacao bool                                     //Se precisa estar logado
}

// Configurar coloca todas as rotas dentro do router
func Configurar(r *mux.Router) *mux.Router {
	rotas := rotasUsuarios                   //Pega as rotas definidas
	rotas = append(rotas, rotaLogin)         //Adiciona a rota de login
	rotas = append(rotas, rotasAdmin...)     //Adiciona as rotas administrativas
	rotas = append(rotas, rotasProdutos...)  //Adiciona as rotas de produtos
	rotas = append(rotas, rotasFavoritos...) //Adiciona as rotas de favoritos

	for _, rota := range rotas { //Para cada rota...

		if rota.RequerAutenticacao {
			r.HandleFunc(rota.URI,
				middlewares.Logger(middlewares.Autenticar(rota.Funcao)),
			).Methods(rota.Metodo) //se a rota requerir autenticacao, chama essa função
		} else {
			r.HandleFunc(rota.URI, middlewares.Logger(rota.Funcao)).Methods(rota.Metodo) //Registra no router: "quando vier requisição X, chame a função Y"
		}

	}

	// Configurar CORS para permitir requisições do frontend
	r.Use(configurarCORS)

	return r //Retorna o router configurado
}

// configurarCORS adiciona os headers CORS necessários
func configurarCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir origem específica (ajuste para seu frontend)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
