package utils

import (
	"os"
	"praktikum3/app/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ambil secret key dari environment (.env)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken membuat JWT token untuk user MongoDB
func GenerateToken(user model.User) (string, error) {
	userID := ""
	if !user.ID.IsZero() {
		userID = user.ID.Hex()
	}

	claims := &model.JWTClaims{
		UserID:   userID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Expired dalam 1 hari
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken memverifikasi JWT dan mengembalikan klaim
func ValidateToken(tokenString string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid {
		// Validasi tambahan: pastikan UserID valid
		if claims.UserID == "" || primitive.IsValidObjectID(claims.UserID) {
			return claims, nil
		}
	}

	return nil, jwt.ErrTokenMalformed
}
