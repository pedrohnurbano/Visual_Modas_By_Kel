package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	//StringConexaoBanco é a string de conexão com o MySQL
	StringConexaoBanco = ""

	//Porta onde a API vai estar rodando
	Porta = 0

	//SecretKey é a chave que vai ser usada para assinar o token
	SecretKey []byte

	//AbacatePayToken é o token de autenticação do AbacatePay
	AbacatePayToken = ""

	//BaseURL é a URL base da aplicação
	BaseURL = ""
)

// Carregar vai inicializar as variáveis de ambiente
func Carregar() {
	var erro error

	if erro := godotenv.Load(); erro != nil { //tudo que tem "env" no nome é uma variável de ambiente que é pega do arquivo .env
		log.Fatal(erro)
	}

	Porta, erro = strconv.Atoi(os.Getenv("API_PORT")) //Atoi pega uma string e converte pra int
	if erro != nil {
		Porta = 9000
	}

	StringConexaoBanco = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USUARIO"),
		os.Getenv("DB_SENHA"),
		os.Getenv("DB_NOME"),
	)

	SecretKey = []byte(os.Getenv("SECRET_KEY"))

	// Carregar token do AbacatePay
	AbacatePayToken = os.Getenv("ABACATEPAY_TOKEN")
	if AbacatePayToken == "" {
		AbacatePayToken = "abc_dev_CtuuPHQhkhGfA0dCAcdnH24D" // Token padrão de desenvolvimento
	}

	// Carregar URL base
	BaseURL = os.Getenv("BASE_URL")
	if BaseURL == "" {
		BaseURL = "http://localhost:3000" // URL padrão para desenvolvimento
	}
}
