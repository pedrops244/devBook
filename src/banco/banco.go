package banco

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var ConexaoAtual *sql.DB // Variável global para a conexão ativa

// Conectar configura uma nova conexão e armazena na variável global
func Conectar(stringConexao string) (*sql.DB, error) {
	db, erro := sql.Open("sqlserver", stringConexao)
	if erro != nil {
		log.Printf("Erro ao abrir a conexão com o banco: %v\n", erro)
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		log.Printf("Erro ao tentar conectar: %v\n", erro)
		return nil, erro
	}

	// Fecha a conexão anterior, se existir
	if ConexaoAtual != nil {
		ConexaoAtual.Close()
	}

	ConexaoAtual = db // Define a nova conexão como a ativa
	return db, nil
}

func ObterConexao() (*sql.DB, error) {
	if ConexaoAtual == nil {
		return nil, sql.ErrConnDone // Retorna erro se não houver conexão ativa
	}
	return ConexaoAtual, nil
}
