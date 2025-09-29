package modelos

// Usuario representa um usuário no sistema
type Usuario struct {
	ID        uint64 `json:"id,omitempty"`
	Nome      string `json:"nome,omitempty"`
	Sobrenome string `json:"sobrenome,omitempty"`
	Email     string `json:"email,omitempty"`
	Telefone  string `json:"telefone,omitempty"`
	CPF       string `json:"cpf,omitempty"`
	Role      string `json:"role,omitempty"`
	CriadoEm  string `json:"criadoEm,omitempty"`
}