package service

import (
	"context"
	"praktikum3/app/model"
	"praktikum3/app/repository"
	"praktikum3/app/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	userRepo repository.IUserRepository
	password utils.PasswordChecker
	tokenGen utils.TokenGenerator
}

// ===============================
// Constructor untuk runtime aplikasi
// ===============================
func NewAuthService(db *mongo.Database) *AuthService {
	repo := repository.NewUserRepository(db)

	return &AuthService{
		userRepo: repo,
		password: utils.RealPasswordChecker{},
		tokenGen: utils.RealTokenGenerator{},
	}
}

// ===============================
// Constructor untuk Unit Test
// ===============================
func NewAuthServiceMock(
	repo repository.IUserRepository,
	pw utils.PasswordChecker,
	tg utils.TokenGenerator,
) *AuthService {
	return &AuthService{
		userRepo: repo,
		password: pw,
		tokenGen: tg,
	}
}

// ========================================
// @Summary Login user
// @Description Login dan mendapatkan JWT token dari sistem
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} map[string]interface{} "Body tidak valid"
// @Failure 401 {object} map[string]interface{} "Username atau password salah"
// @Failure 500 {object} map[string]interface{} "Kesalahan server atau database"
// @Router /login [post]
// ========================================
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Body tidak valid",
		})
	}

	// ambil user dari database
	user, err := s.userRepo.FindByUsernameOrEmail(context.Background(), req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Kesalahan database: " + err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Username atau password salah",
		})
	}

	// verify password menggunakan dependency injection
	if !s.password.Check(user.PasswordHash, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Username atau password salah",
		})
	}

	// generate JWT token menggunakan dependency injection
	token, err := s.tokenGen.Generate(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat token",
		})
	}

	// response sukses
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
