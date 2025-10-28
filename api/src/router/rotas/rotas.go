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
	rotas = append(rotas, rotasCarrinho...)  //Adiciona as rotas de carrinho
	rotas = append(rotas, rotasPedidos...)   //Adiciona as rotas de pedidos

	for _, rota := range rotas { //Para cada rota...

		if rota.RequerAutenticacao {
			r.HandleFunc(rota.URI,
				middlewares.Logger(middlewares.Autenticar(rota.Funcao)),
			).Methods(rota.Metodo) //se a rota requerir autenticacao, chama essa função
		} else {
			r.HandleFunc(rota.URI, middlewares.Logger(rota.Funcao)).Methods(rota.Metodo) //Registra no router: "quando vier requisição X, chame a função Y"
		}

	}

	return r //Retorna o router configurado
}
