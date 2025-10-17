package repository

import (
	"database/sql"
	"praktikum3/app/model"
	"time"
)

type UserRepository struct {
	DB *sql.DB
}

// Constructor
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Ambil user berdasarkan username atau email (untuk login)
func (r *UserRepository) FindByUsernameOrEmail(usernameOrEmail string) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, role, created_at
		FROM users
		WHERE username = $1 OR email = $1
	`, usernameOrEmail).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Soft delete user berdasarkan ID
func (r *UserRepository) SoftDeleteUser(id uint) error {
	query := `UPDATE users SET deleted_at = $1 WHERE id = $2`
	_, err := r.DB.Exec(query, time.Now(), id)
	return err
}
