package utils

import (
	"praktikum3/app/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Pakai secret yang sama persis di middleware
var jwtSecret = []byte("rahasia-super") // jangan beda dengan middleware

// GenerateToken membuat JWT token untuk user
func GenerateToken(user model.User) (string, error) {
	claims := &model.JWTClaims{ // pointer lebih aman
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // expired 1 hari
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken memverifikasi JWT
func ValidateToken(tokenString string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// langsung return claims
	if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenMalformed
}
