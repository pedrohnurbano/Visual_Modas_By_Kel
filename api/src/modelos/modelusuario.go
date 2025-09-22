package modelos

import (
	"api/src/seguranca"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// Usuario representa um usuário utilizando a rede social
type Usuario struct {
	ID        uint64    `json:"id,omitempty"` //omitempty faz com que o campo seja ignorado se estiver vazio
	Nome      string    `json:"nome,omitempty"`
	Sobrenome string    `json:"sobrenome,omitempty"`
	Email     string    `json:"email,omitempty"`
	Senha     string    `json:"senha,omitempty"`
	Telefone  string    `json:"telefone,omitempty"`
	CPF       string    `json:"cpf,omitempty"`
	Role      string    `json:"role,omitempty"` // Nova campo para role (user ou admin)
	CriadoEm  time.Time `json:"CriadoEm,omitempty"`
}

// Preparar vai chamar os métodos para validar e formatar o usuário recebido
func (usuario *Usuario) Preparar(etapa string) error {
	if erro := usuario.validar(etapa); erro != nil {
		return erro
	}

	if erro := usuario.formatar(etapa); erro != nil {
		return erro
	}

	return nil
}

func (usuario *Usuario) validar(etapa string) error {
	if usuario.Nome == "" {
		return errors.New("o nome é obrigatório e não pode estar em branco")
	}

	if usuario.Sobrenome == "" {
		return errors.New("o sobrenome é obrigatório e não pode estar em branco")
	}

	if usuario.Email == "" {
		return errors.New("o e-mail é obrigatório e não pode estar em branco")
	}

	if usuario.Telefone == "" {
		return errors.New("o telefone é obrigatório e não pode estar em branco")
	}

	if usuario.CPF == "" {
		return errors.New("o CPF é obrigatório e não pode estar em branco")
	}

	if erro := checkmail.ValidateFormat(usuario.Email); erro != nil { //valida o formato do e-mail
		return errors.New("o e-mail inserido é inválido") //se o e-mail for inválido, retorna um erro
	}

	if etapa == "cadastro" && usuario.Senha == "" {
		return errors.New("a senha é obrigatória e não pode estar em branco")
	}

	// Validação da role
	if usuario.Role != "" && usuario.Role != "user" && usuario.Role != "admin" {
		return errors.New("role inválida. Deve ser 'user' ou 'admin'")
	}

	return nil
}

func (usuario *Usuario) formatar(etapa string) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome) //Remove espaços em branco no início e no final
	usuario.Sobrenome = strings.TrimSpace(usuario.Sobrenome)
	usuario.Email = strings.TrimSpace(usuario.Email)
	usuario.Telefone = strings.TrimSpace(usuario.Telefone)
	usuario.CPF = strings.TrimSpace(usuario.CPF)

	// Se a role não for especificada ou for vazia, define como "user" por padrão
	if usuario.Role == "" {
		usuario.Role = "user"
	}

	if etapa == "cadastro" {
		senhaComHash, erro := seguranca.Hash(usuario.Senha) //Hasheia a senha
		if erro != nil {
			return erro
		}

		usuario.Senha = string(senhaComHash) //Converte o hash de volta para string
	}

	return nil
}
