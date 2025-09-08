package main

import (
	"fmt"
	"log"
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/config"
	"visual_modas_by_kel/visual_modas_by_kel/src/router"
	"visual_modas_by_kel/visual_modas_by_kel/src/utils"
)

func main() {
	config.Carregar()
	utils.CarregarTemplates()
	r := router.Gerar()

	fmt.Printf("Escutando na porta %d", config.Porta)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))
}
