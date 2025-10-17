package service

import (
	"database/sql"
	"praktikum3/app/model"
	"praktikum3/app/repository"
	"praktikum3/app/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

// Constructor
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{
		UserRepo: repository.NewUserRepository(db),
	}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Body tidak valid",
		})
	}

	user, err := s.UserRepo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Username atau password salah",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Kesalahan database",
		})
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Username atau password salah",
		})
	}

	token, err := utils.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat token",
		})
	}

	return c.JSON(model.LoginResponse{
		User: model.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	})
}
