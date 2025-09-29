package middlewares

import (
	"log"
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/cookies"
)

// Logger escreve informações da requisição no terminal
func Logger(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %s %s %s", r.Method, r.RequestURI, r.Host)
		proximaFuncao(w, r)
	}
}

// Autenticar verifica a existência do cookie de autenticação
func Autenticar(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, erro := cookies.Ler(r); erro != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		proximaFuncao(w, r)
	}
}

// AutenticarAdmin verifica se o usuário é admin
func AutenticarAdmin(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, erro := cookies.Ler(r)
		if erro != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Verificar se o usuário tem role de admin
		if cookie["role"] != "admin" {
			http.Redirect(w, r, "/home", http.StatusForbidden)
			return
		}

		proximaFuncao(w, r)
	}
}
