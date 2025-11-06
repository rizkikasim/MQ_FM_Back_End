package utils

import (
	"time"

	"mqfm_backend/internal/helper"

)

// Wrapper untuk JWT biar service gak import helper langsung
func GenerateJWT(email, secretKey string, duration time.Duration) (string, error) {
	return helper.GenerateJWT(email, secretKey, duration)
}

func ParseJWT(tokenString, secretKey string) (string, error) {
	return helper.ParseJWT(tokenString, secretKey)
}
