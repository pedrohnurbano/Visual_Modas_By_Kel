package main

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	config.Carregar()
	
	// Criar pasta de uploads se não existir
	if erro := os.MkdirAll("./uploads/produtos", 0755); erro != nil {
		log.Printf("Aviso: Não foi possível criar pasta de uploads: %v", erro)
	}
	
	r := router.Gerar()

	// Servir arquivos estáticos da pasta uploads
	fs := http.FileServer(http.Dir("./uploads"))
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))

	fmt.Printf("Escutando na porta %d\n", config.Porta)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))
}
