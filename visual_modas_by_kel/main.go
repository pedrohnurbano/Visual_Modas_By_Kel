package main

import (
	"fmt"
	"log"
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/router"
	"visual_modas_by_kel/visual_modas_by_kel/src/utils"
)

func main() {
	utils.CarregarTemplates()
	r := router.Gerar()

	fmt.Println("Escutando na porta 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
