package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"praktikum3/app/model"
	"praktikum3/app/service"
	"praktikum3/tests/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestApp(authSvc *service.AuthService) *fiber.App {
	app := fiber.New()
	app.Post("/login", authSvc.Login)
	return app
}

// ===========================
// 1. INVALID BODY
// ===========================
func TestLogin_InvalidBody(t *testing.T) {
	mockRepo := &mocks.UserRepositoryMock{}
	mockPw := mocks.PasswordCheckerMock{}
	mockToken := mocks.TokenGeneratorMock{}

	authSvc := service.NewAuthServiceMock(mockRepo, mockPw, mockToken)
	app := setupTestApp(authSvc)

	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("NOT_JSON"))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

// ===========================
// 2. USER NOT FOUND
// ===========================
func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := &mocks.UserRepositoryMock{
		FindByUsernameOrEmailFunc: func(ctx context.Context, username string) (*model.User, error) {
			return nil, nil
		},
	}

	mockPw := mocks.PasswordCheckerMock{}
	mockToken := mocks.TokenGeneratorMock{}

	authSvc := service.NewAuthServiceMock(mockRepo, mockPw, mockToken)
	app := setupTestApp(authSvc)

	body, _ := json.Marshal(model.LoginRequest{Username: "user", Password: "123"})
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}

// ===========================
// 3. PASSWORD SALAH
// ===========================
func TestLogin_WrongPassword(t *testing.T) {
	mockUser := &model.User{
		ID:           primitive.NewObjectID(),
		Username:     "user",
		Email:        "u@mail.com",
		PasswordHash: "HASHED",
		Role:         "alumni",
		CreatedAt:    time.Now(),
	}

	mockRepo := &mocks.UserRepositoryMock{
		FindByUsernameOrEmailFunc: func(ctx context.Context, username string) (*model.User, error) {
			return mockUser, nil
		},
	}

	mockPw := mocks.PasswordCheckerMock{
		CheckFunc: func(hash, password string) bool { return false },
	}

	mockToken := mocks.TokenGeneratorMock{}

	authSvc := service.NewAuthServiceMock(mockRepo, mockPw, mockToken)
	app := setupTestApp(authSvc)

	body, _ := json.Marshal(model.LoginRequest{Username: "user", Password: "wrong"})
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}

// ===========================
// 4. TOKEN GENERATOR ERROR
// ===========================
func TestLogin_TokenError(t *testing.T) {
	mockUser := &model.User{
		ID:           primitive.NewObjectID(),
		Username:     "user",
		Email:        "u@mail.com",
		PasswordHash: "HASHED",
		Role:         "alumni",
		CreatedAt:    time.Now(),
	}

	mockRepo := &mocks.UserRepositoryMock{
		FindByUsernameOrEmailFunc: func(ctx context.Context, username string) (*model.User, error) {
			return mockUser, nil
		},
	}

	mockPw := mocks.PasswordCheckerMock{
		CheckFunc: func(hash, password string) bool { return true },
	}

	mockToken := mocks.TokenGeneratorMock{
		GenerateFunc: func(user model.User) (string, error) {
			return "", errors.New("token error")
		},
	}

	authSvc := service.NewAuthServiceMock(mockRepo, mockPw, mockToken)
	app := setupTestApp(authSvc)

	body, _ := json.Marshal(model.LoginRequest{Username: "user", Password: "123"})
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

// ===========================
// 5. LOGIN SUKSES
// ===========================
func TestLogin_Success(t *testing.T) {
	mockUser := &model.User{
		ID:           primitive.NewObjectID(),
		Username:     "user",
		Email:        "u@mail.com",
		PasswordHash: "HASHED",
		Role:         "alumni",
		CreatedAt:    time.Now(),
	}

	mockRepo := &mocks.UserRepositoryMock{
		FindByUsernameOrEmailFunc: func(ctx context.Context, username string) (*model.User, error) {
			return mockUser, nil
		},
	}

	mockPw := mocks.PasswordCheckerMock{
		CheckFunc: func(hash, password string) bool { return true },
	}

	mockToken := mocks.TokenGeneratorMock{
		GenerateFunc: func(user model.User) (string, error) {
			return "TOKEN123", nil
		},
	}

	authSvc := service.NewAuthServiceMock(mockRepo, mockPw, mockToken)
	app := setupTestApp(authSvc)

	body, _ := json.Marshal(model.LoginRequest{Username: "user", Password: "123"})
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}
