package respostas

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON vai receber um statusCode, colocar ele no header, e pegar os dados genéricos e transformar em JSON
func JSON(w http.ResponseWriter, statusCode int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json") //setando o tipo de conteúdo a ser retornado como JSON
	w.WriteHeader(statusCode)                          //WriteHeader retorna o código de status

	if dados != nil { //se dados não for nulo, ele vai transformar em JSON
		if erro := json.NewEncoder(w).Encode(dados); erro != nil { //Encoder transforma dados em JSON
			log.Fatal(erro)
		}
	}
}

// Erro rtetorna um erro em formato JSON
func Erro(w http.ResponseWriter, statusCode int, erro error) {
	JSON(w, statusCode, struct {
		Erro string `json:"erro"`
	}{
		Erro: erro.Error(),
	})
}
