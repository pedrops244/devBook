package banco

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

var pool *sql.DB // Variável global para o pool de conexões

// Conectar configura um pool de conexões com os parâmetros fornecidos.
func Conectar(stringConexao string) error {
	if pool != nil {
		pool.Close() // Fecha o pool atual antes de configurar um novo
	}

	db, err := sql.Open("sqlserver", stringConexao)
	if err != nil {
		log.Printf("Erro ao abrir a conexão com o banco: %v\n", err)
		return err
	}

	// Configura limites do pool
	db.SetMaxOpenConns(10)                  // Máximo de conexões abertas
	db.SetMaxIdleConns(5)                   // Máximo de conexões ociosas
	db.SetConnMaxLifetime(30 * time.Minute) // Tempo máximo de vida das conexões

	if err = db.Ping(); err != nil {
		db.Close()
		log.Printf("Erro ao tentar conectar: %v\n", err)
		return err
	}

	pool = db
	log.Println("Pool de conexões configurado com sucesso")
	return nil
}

// ObterConexao retorna uma conexão do pool.
func ObterConexao() (*sql.DB, error) {
	if pool == nil {
		return nil, sql.ErrConnDone // Erro se o pool ainda não foi configurado
	}
	return pool, nil
}
