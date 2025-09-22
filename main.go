package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Produto struct {
	ID        int     `json:"id"`
	Nome      string  `json:"nome"`
	Categoria string  `json:"categoria"`
	Preco     float64 `json:"preco"`
	Estoque   int     `json:"estoque"`
	Status    string  `json:"status"`
	SKU       string  `json:"sku"`
}

func main() {
	r := mux.NewRouter()

	// Serve arquivos estáticos (HTML, CSS, JS, imagens)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./visual_modas_by_kel/")))

	// Endpoint de exemplo para produtos (GET)
	r.HandleFunc("/api/produtos", ProdutosHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Servidor rodando em http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func ProdutosHandler(w http.ResponseWriter, r *http.Request) {
	// Exemplo estático, depois trocar para buscar do MySQL
	produtos := []Produto{
		{ID: 1, Nome: "Produto Exemplo", Categoria: "Categoria", Preco: 99.90, Estoque: 10, Status: "ativo", SKU: "SKU001"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produtos)
}
