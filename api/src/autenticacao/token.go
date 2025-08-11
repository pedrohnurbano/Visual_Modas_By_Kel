package autenticacao

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//CriarToken retorna um token assinado com as permissões do usuário
func CriarToken(usuarioID uint64) (string, error) { //retorna o token e um erro
	permissoes := jwt.MapClaims{} //map que vai ter as permissoes dentro do token
	permissoes["authorized"] = true //quem tiver esse campo está autorizado
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix() //token expira depois de 6h, .unix converte para a anotação unix de data, devolve a quantidade de segundos que passaram desde o dia 01/01/1970
	permissoes["usuarioId"] = usuarioID //usuario que é dono do token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes) //criando token novo a partir dos claims que foram criados, passa o método de signing HS256, e as permissoes que ele vai ter
	return token.SignedString([]byte(config.SecretKey)) //secretkey, chave que vai ser usada para fazer a assinatura e garantir a autenticidade do token
}

//ValidarToken verifica se o token passado na requisição é válido
func ValidarToken(r *http.Request) error {
	tokenString := extrairToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if erro != nil {
		return erro
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("token inválido")
}

//ExtrairUsuarioID retorna o usuarioId que está salvo no token
func ExtrairUsuarioID(r *http.Request) (uint64, error) {
	tokenString := extrairToken(r)
	token, erro := jwt.Parse(tokenString, retornarChaveDeVerificacao)
	if erro != nil {
		return 0, erro
	}

	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usuarioID, erro := strconv.ParseUint(fmt.Sprintf("%.0f", permissoes["usuarioId"]), 10, 64) //converte o id do usuario pra uint64
		if erro != nil {
			return 0, erro
		}

		return usuarioID, nil
	}

	return 0, errors.New("token invalido")
}

func extrairToken(r *http.Request) string { //pega oq tiver no authorization do request e joga na variavel token
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 { //se tiver o tamanho de 2, retorna a 2° palavra
		return strings.Split(token, " ")[1] //2° palavra
	}

	return "" //se n tiver duas posicoes retorna uma string em branco
}

func retornarChaveDeVerificacao(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("método de assinatura inesperado! %v", token.Header["alg"])
	}

	return []byte(config.SecretKey), nil
}
