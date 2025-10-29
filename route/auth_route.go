package route

import (
	"praktikum3/app/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthRoute mendaftarkan semua endpoint autentikasi
func AuthRoute(app fiber.Router, db *mongo.Database) {
	authService := service.NewAuthService(db)

	// 🟢 Endpoint login (tanpa middleware)
	app.Post("/login", authService.Login)
}
