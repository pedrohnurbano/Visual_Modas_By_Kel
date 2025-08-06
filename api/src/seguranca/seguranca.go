package seguranca

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash recebe uma string e coloca um hash nela
func Hash(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost) //Gera o hash da senha com o custo padrão
}

//VerificarSenha compara uma senha e um hash e retorna se elas são iguais
func VerificarSenha(SenhaComHash, senhaString string) error {
	return bcrypt.CompareHashAndPassword([]byte(SenhaComHash), []byte(senhaString)) //Compara a senha com o hash da senha
}

//HASHS:
//Cadastro - A senha é hasheada e só o hash vai pro banco (Hash)
//Login - A senha digitada é comparada com o hash guardado (VerificarSenha)
//Segurança - Mesmo que alguém roube o banco, não consegue ver as senhas
