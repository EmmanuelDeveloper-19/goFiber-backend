package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "secreto_super_seguro"

func GenerateJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":   userID,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// GetSecretKey retorna la clave secreta usada para firmar tokens JWT
func GetSecretKey() string {
	return secretKey
}
