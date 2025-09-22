package modelos

// DadosAutenticacao contém o id, token e role do usuário autenticado
type DadosAutenticacao struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Role  string `json:"role"`
}
