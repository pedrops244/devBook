package config

import (
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

	// Carrega o arquivo .env
	if erro = godotenv.Load(); erro != nil {
		log.Fatal("Erro ao carregar o arquivo .env", erro)
	}

	// Carrega a porta da API
	Porta, erro = strconv.Atoi(os.Getenv("API_PORT"))
	if erro != nil {
		log.Fatal("API_PORT não configurado corretamente no .env")
	}

	// Carrega a chave secreta do JWT
	SecretKey = []byte(os.Getenv("JWT_SECRET"))
	if len(SecretKey) == 0 {
		log.Fatal("JWT_SECRET não configurado corretamente no .env")
	}
}
