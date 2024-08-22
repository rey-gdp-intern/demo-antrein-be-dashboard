package entity

import "github.com/golang-jwt/jwt/v5"

type JWTClaim struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
