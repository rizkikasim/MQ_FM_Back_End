package helper

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"

)

// ParseJWT men-decode token JWT dan mengembalikan email dari claims
func ParseJWT(tokenString, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing tidak valid")
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("token tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("gagal membaca claims JWT")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("email tidak ditemukan di token")
	}

	return email, nil
}
