package controllers

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/utils"
)

// CarregarPaginaFavoritos carrega a página de favoritos
func CarregarPaginaFavoritos(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "favoritos.html", nil)
}

// CarregarPaginaSacola carrega a página da sacola
func CarregarPaginaSacola(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "sacola.html", nil)
}

// CarregarPaginaRoupas carrega a página de listagem de todas as roupas
func CarregarPaginaRoupas(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "roupas.html", nil)
}

// CarregarPaginaCheckout carrega a página de checkout
func CarregarPaginaCheckout(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "checkout.html", nil)
}

// CarregarPaginaTermosPolitica carrega a página de termos e política
func CarregarPaginaTermosPolitica(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "termos_politica.html", nil)
}
