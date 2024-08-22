package generator

import (
	"antrein/bc-dashboard/model/entity"
	"math/rand"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(key string, claims entity.JWTClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func GenerateRandomString(lenStr int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, lenStr)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
