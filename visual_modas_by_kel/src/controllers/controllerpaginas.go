package controllers

import (
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/utils"
)

// CarregarPaginaFavoritos carrega a página de favoritos
func CarregarPaginaFavoritos(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "favoritos.html", nil)
}
