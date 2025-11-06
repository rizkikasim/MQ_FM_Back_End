package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

)

// GenerateJWT membuat token JWT dengan payload email dan masa berlaku tertentu
func GenerateJWT(email, secretKey string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
