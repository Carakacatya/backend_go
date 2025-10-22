package route

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"praktikum3/app/service"
)

func AuthRoute(app fiber.Router, db *sql.DB) {
	authService := service.NewAuthService(db)

	// Endpoint login (tidak perlu middleware)
	app.Post("/login", authService.Login)
}
