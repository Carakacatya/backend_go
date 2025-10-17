package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	ID           int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string         `json:"username" gorm:"unique;not null"`
	Email        string         `json:"email" gorm:"unique;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	Role         string         `json:"role" gorm:"type:enum('admin','user');default:'user'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
