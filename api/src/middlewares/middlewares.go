package middlewares

import (
	"api/src/autenticacao"
	"api/src/respostas"
	"log"
	"net/http"
)
//middleware -> camada entre a requisição e a resposta
//Logger escreve informações da requisição no terminal
func Logger(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %s %s %s", r.Method, r.RequestURI, r.Host) //exibindo informações da requisição no terminal
		proximaFuncao(w, r)
	}
}

//Autenticar verifica se o usuário fazendo a requisição está autenticando
func Autenticar(proximaFuncao http.HandlerFunc) http.HandlerFunc { //handlerfunc recebe um w responsewriter e um r *request
	return func(w http.ResponseWriter, r *http.Request) { //vai chamar a função que vai validar o token, e vai executar o q veio no parametro (proximaFuncao)
		if erro := autenticacao.ValidarToken(r); erro != nil {
			respostas.Erro(w, http.StatusUnauthorized, erro)
			return
		}
		proximaFuncao(w, r)
	}
}
