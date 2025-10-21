package main

import (
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
)

// Script auxiliar para gerar o hash de uma senha

func main() {
	
	senhaAdmin := "123456" // Senha que vai ser hasheada
	
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
