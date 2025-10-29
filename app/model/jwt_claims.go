package model

import "github.com/golang-jwt/jwt/v5"

// JWTClaims menyimpan payload di dalam token JWT
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
