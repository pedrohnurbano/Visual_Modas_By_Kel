package requisicoes

import (
	"bytes"
	"fmt"
	"net/http"
	"visual_modas_by_kel/visual_modas_by_kel/src/cookies"
)

// FazerRequisicaoComAutenticacao adiciona o token de autenticação na requisição
func FazerRequisicaoComAutenticacao(r *http.Request, metodo, url string, dados []byte) (*http.Response, error) {
	request, erro := http.NewRequest(metodo, url, bytes.NewBuffer(dados))
	if erro != nil {
		return nil, erro
	}

	cookie, erro := cookies.Ler(r)
	if erro != nil {
		return nil, erro
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cookie["token"]))
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, erro := client.Do(request)
	if erro != nil {
		return nil, erro
	}

	return response, nil
}