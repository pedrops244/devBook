package auth

import (
	"api/src/config"
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
