package service

import (
	"context"
	"fmt"
	"praktikum3/app/model"
	"praktikum3/app/repository"
	"praktikum3/app/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	userRepo repository.IUserRepository
}

// ✅ Constructor
func NewAuthService(db *mongo.Database) *AuthService {
	repo := repository.NewUserRepository(db)
	return &AuthService{
		userRepo: repo,
	}
}

// ✅ Login Handler
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Body tidak valid",
		})
	}

	// 🔍 Cari user berdasarkan username/email dari MongoDB
	user, err := s.userRepo.FindByUsernameOrEmail(context.Background(), req.Username)
	fmt.Println("DEBUG user =", user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Kesalahan database: " + err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Username atau password salah (user tidak ditemukan)",
		})
	}

	// 🔑 Verifikasi password
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		fmt.Println("❌ Password tidak cocok antara:", req.Password, "dan hash:", user.PasswordHash)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Username atau password salah",
		})
	}

	// 🪪 Generate JWT token
	token, err := utils.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat token: " + err.Error(),
		})
	}

	// ✅ Respons sukses
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
