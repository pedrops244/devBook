package config

import "fmt"

type DBConfig struct {
	UsuarioDB string `json:"usuarioDB"`
	SenhaDB   string `json:"senhaDB"`
	Host      string `json:"host"`
	Database  string `json:"database"`
}

func ConstruirStringConexao(dbConfig DBConfig) string {
	return fmt.Sprintf(
		"sqlserver://%s:%s@%s?database=%s&encrypt=disable&trusted_connection=true&tcpKeepAlive=1",
		dbConfig.UsuarioDB, dbConfig.SenhaDB, dbConfig.Host, dbConfig.Database,
	)
}
