package repositorios

import (
	"api/src/modelos"
	"database/sql"
	"fmt"
	"log"
)

// Usuarios representa um repositório de usuarios
type Usuarios struct {
	db *sql.DB
}

// NovoRepositorioDeUsuarios é uma função que recebe o banco de dados e retorna um novo repositorio de usuarios
func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

// Criar insere um usuario no banco de dados
func (repositorio Usuarios) Criar(usuario modelos.Usuario) (uint64, error) {
	log.Printf("Tentando criar usuário: %+v", usuario)
	
	// Garante que novos usuários sempre sejam criados como "user"
	if usuario.Role == "" {
		usuario.Role = "user"
	}
	
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO usuarios (nome, sobrenome, email, senha, telefone, cpf, role) VALUES(?, ?, ?, ?, ?, ?, ?)",
	)
	if erro != nil {
		log.Printf("Erro ao preparar statement: %v", erro)
		return 0, erro
	}
	defer statement.Close()

	log.Printf("Statement preparado. Executando com valores: nome=%s, sobrenome=%s, email=%s, telefone=%s, cpf=%s, role=%s", 
		usuario.Nome, usuario.Sobrenome, usuario.Email, usuario.Telefone, usuario.CPF, usuario.Role)

	resultado, erro := statement.Exec(usuario.Nome, usuario.Sobrenome, usuario.Email, usuario.Senha, usuario.Telefone, usuario.CPF, usuario.Role)
	if erro != nil {
		log.Printf("Erro ao executar statement: %v", erro)
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		log.Printf("Erro ao obter último ID: %v", erro)
		return 0, erro
	}

	log.Printf("Usuário inserido com sucesso. ID: %d, Role: %s", ultimoIDInserido, usuario.Role)
	return uint64(ultimoIDInserido), nil
}

// Buscar traz todos os usuarios que atendem um filtro de nome, sobrenome ou email
func (repositorio Usuarios) Buscar(nomeOuEmail string) ([]modelos.Usuario, error) {
	nomeOuEmail = fmt.Sprintf("%%%s%%", nomeOuEmail)
	
	linhas, erro := repositorio.db.Query(
		"SELECT id, nome, sobrenome, email, telefone, cpf, role, criadoEm FROM usuarios WHERE nome LIKE ? OR sobrenome LIKE ? OR email LIKE ?",
		nomeOuEmail, nomeOuEmail, nomeOuEmail,
	)

	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []modelos.Usuario

	for linhas.Next() {
		var usuario modelos.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Sobrenome,
			&usuario.Email,
			&usuario.Telefone,
			&usuario.CPF,
			&usuario.Role,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

// BuscarPorID traz um usuario pelo seu id
func (repositorio Usuarios) BuscarPorID(ID uint64) (modelos.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"SELECT id, nome, sobrenome, email, telefone, cpf, role, criadoEm FROM usuarios WHERE id = ?",
		ID,
	)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linhas.Close()

	var usuario modelos.Usuario

	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Sobrenome,
			&usuario.Email,
			&usuario.Telefone,
			&usuario.CPF,
			&usuario.Role,
			&usuario.CriadoEm,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

// Atualizar altera as informações de um usuario no banco de dados
func (repositorio Usuarios) Atualizar(ID uint64, usuario modelos.Usuario) error {
	// Não permite atualizar a role através desta função
	statement, erro := repositorio.db.Prepare(
		"UPDATE usuarios SET nome = ?, sobrenome = ?, email = ?, telefone = ?, cpf = ? WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuario.Nome, usuario.Sobrenome, usuario.Email, usuario.Telefone, usuario.CPF, ID); erro != nil {
		return erro
	}

	return nil
}

// Deletar exclui as informações de um usuário no banco de dados
func (repositorio Usuarios) Deletar(ID uint64) error {
	statement, erro := repositorio.db.Prepare("DELETE FROM usuarios WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}

// BuscarPorEmail busca um usuário por email e retorna o seu id, senha e role
func (repositorio Usuarios) BuscarPorEmail(email string) (modelos.Usuario, error) {
	linha, erro := repositorio.db.Query("SELECT id, senha, role FROM usuarios WHERE email = ?", email)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linha.Close()

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Senha, &usuario.Role); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

// BuscarPorCPF busca um usuário por CPF (útil para verificar duplicatas)
func (repositorio Usuarios) BuscarPorCPF(cpf string) (modelos.Usuario, error) {
	linha, erro := repositorio.db.Query("SELECT id, nome, sobrenome, email, telefone, cpf, role, criadoEm FROM usuarios WHERE cpf = ?", cpf)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linha.Close()

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Sobrenome,
			&usuario.Email,
			&usuario.Telefone,
			&usuario.CPF,
			&usuario.Role,
			&usuario.CriadoEm,
		); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

func (repositorio Usuarios) BuscarSenha(usuarioID uint64) (string, error) {
	linha, erro := repositorio.db.Query("SELECT senha FROM usuarios WHERE id = ?", usuarioID)
	if erro != nil {
		return "", erro
	}
	defer linha.Close()

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.Senha); erro != nil {
			return "", erro
		}
	}

	return usuario.Senha, nil
}

// AtualizarSenha altera a senha de um usuário
func (repositorio Usuarios) AtualizarSenha(usuarioID uint64, senha string) error {
	statement, erro := repositorio.db.Prepare("UPDATE usuarios SET senha = ? WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(senha, usuarioID); erro != nil {
		return erro
	}

	return nil
}

// AtualizarRole altera a role de um usuário (apenas para admins)
func (repositorio Usuarios) AtualizarRole(usuarioID uint64, role string) error {
	statement, erro := repositorio.db.Prepare("UPDATE usuarios SET role = ? WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(role, usuarioID); erro != nil {
		return erro
	}

	return nil
}
