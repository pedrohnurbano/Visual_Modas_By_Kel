package config

import (
	"fmt"
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

	// Tentar carregar .env, mas não falhar se não existir
	if erro = godotenv.Load(); erro != nil {
		log.Printf("Arquivo .env não encontrado, usando valores padrão: %v", erro)
	}

	// Usar valores padrão se as variáveis de ambiente não estiverem definidas
	if os.Getenv("APP_PORT") != "" {
		Porta, erro = strconv.Atoi(os.Getenv("APP_PORT"))
		if erro != nil {
			log.Fatalf("Erro ao converter a porta da aplicação: %v", erro)
		}
	} else {
		Porta = 3000 // Valor padrão
	}

	if os.Getenv("API_URL") != "" {
		APIURL = os.Getenv("API_URL")
	}
	// APIURL já tem valor padrão definido acima

	if os.Getenv("HASH_KEY") != "" {
		HashKey = []byte(os.Getenv("HASH_KEY"))
	} else {
		HashKey = []byte("minha-chave-hash-super-secreta")
	}

	if os.Getenv("BLOCK_KEY") != "" {
		BlockKey = []byte(os.Getenv("BLOCK_KEY"))
	} else {
		BlockKey = []byte("minha-chave-block-super-secreta-32-chars")
	}

	fmt.Printf("Configuração carregada - Porta: %d, API URL: %s\n", Porta, APIURL)
}
