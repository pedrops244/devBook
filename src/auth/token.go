package auth

import (
	"api/src/config"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CriarToken(usuarioId uint) (string, error) {
	permissioes := jwt.MapClaims{}
	permissioes["authorized"] = true
	permissioes["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissioes["usuarioId"] = usuarioId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissioes)
	return token.SignedString([]byte(config.SecretKey))
}

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

func extrairToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func retornarChaveDeVerificacao(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("método de assinatura inesperado %v", token.Header["alg"])
	}

	return config.SecretKey, nil

}
