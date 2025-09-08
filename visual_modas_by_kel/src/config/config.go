package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	//ApiUrl representa a URL para comunicação com a API
	APIURL = "http://localhost:5000"
	//Porta onde a aplicação web estará rodando
	Porta = 0
	//HashKey é utilizada para autenticação dos cookies
	HashKey []byte
	//BlockKey é utilizada para criptografia dos cookies
	BlockKey []byte
)

// Carregar inicializa as variáveis de ambiente
func Carregar() {
	var erro error

	if erro = godotenv.Load(); erro != nil {
		log.Fatalf("Erro ao carregar as variáveis de ambiente: %v", erro)
	}

	Porta, erro = strconv.Atoi(os.Getenv("APP_PORT"))
	if erro != nil {
		log.Fatalf("Erro ao converter a porta da aplicação: %v", erro)
	}

	APIURL = os.Getenv("API_URL")
	HashKey = []byte(os.Getenv("HASH_KEY"))
	BlockKey = []byte(os.Getenv("BLOCK_KEY"))

}
