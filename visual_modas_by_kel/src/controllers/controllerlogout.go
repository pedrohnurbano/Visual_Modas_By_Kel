package controllers

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/cookies"
)

// FazerLogout remove os dados de autenticação salvos no browser do usuário
func FazerLogout(w http.ResponseWriter, r *http.Request) {
	cookies.Deletar(w)
	http.Redirect(w, r, "/login", http.StatusFound)
}
