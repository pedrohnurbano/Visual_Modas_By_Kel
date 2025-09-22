package main

import (
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
)

// Este é um script auxiliar para gerar o hash de uma senha
// Use-o para criar a senha do administrador inicial no banco de dados
func main() {
	// MUDE ESTA SENHA PARA A SENHA DESEJADA DO ADMIN
	senhaAdmin := "123456" // <-- ALTERE AQUI
	
	senhaComHash, erro := bcrypt.GenerateFromPassword([]byte(senhaAdmin), bcrypt.DefaultCost)
	if erro != nil {
		log.Fatal(erro)
	}
	
	fmt.Println("Senha original:", senhaAdmin)
	fmt.Println("Hash gerado:", string(senhaComHash))
	fmt.Println("\nUse este hash no script SQL para inserir o administrador:")
	fmt.Printf("\nINSERT INTO usuarios (nome, sobrenome, email, senha, telefone, cpf, role) \n")
	fmt.Printf("VALUES ('Admin', 'Sistema', 'admin@visualmodasbykel.com', '%s', '00000000000', '00000000000', 'admin');\n", string(senhaComHash))
}
