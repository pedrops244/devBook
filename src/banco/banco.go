package banco

import (
	"api/src/config"
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func Conectar() (*sql.DB, error) {
	db, erro := sql.Open("sqlserver", config.StringConexaoBanco)
	if erro != nil {
		log.Printf("Erro ao abrir a conexão com o banco: %v\n", erro)
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		log.Printf("Erro ao tentar conectar (Ping): %v\n", erro)
		return nil, erro
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso!")
	return db, nil
}
