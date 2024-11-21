package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	StringConexaoBanco = ""
	Porta              = 0
	SecretKey          []byte
)

func Carregar() {
	var erro error

	if erro = godotenv.Load(); erro != nil {
		log.Fatal(erro)
	}

	Porta, erro = strconv.Atoi(os.Getenv("API_PORT"))
	if erro != nil {
		Porta = 9000
	}

	// StringConexaoBanco = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
	// 	os.Getenv("DB_USUARIO"),
	// 	os.Getenv("DB_SENHA"),
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_NOME"),
	// )

	StringConexaoBanco = fmt.Sprintf("sqlserver://%s?database=%s&encrypt=disable&trusted_connection=true", os.Getenv("DB_HOST"), os.Getenv("DB_NOME"))

	SecretKey = []byte(os.Getenv("JWT_SECRET"))
}
