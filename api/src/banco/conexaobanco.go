package banco

import (
	"api/src/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //Driver
)

//Conectar abre a conexão com o banco de dados e a retorna
func Conectar() (*sql.DB, error) {
	db, erro := sql.Open("mysql", config.StringConexaoBanco) //Open abre a conexão
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil { //Ping verifica se a conexão está aberta
		db.Close()
		return nil, erro
	}

	return db, nil
	
}
