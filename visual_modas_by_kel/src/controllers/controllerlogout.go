package controllers

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/cookies"
)

// FazerLogout remove os dados de autenticação salvos no browser do usuário
func FazerLogout(w http.ResponseWriter, r *http.Request) {
	cookies.Deletar(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"mensagem": "Logout realizado com sucesso"}`))
}
