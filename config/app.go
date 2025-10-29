package config

import (
	"encoding/json"
	"praktikum3/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewApp(db *mongo.Database) *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ")
		},
		JSONDecoder: json.Unmarshal,
	})

	// === Global middleware ===
	app.Use(middleware.LoggerMiddleware)

	// === Root endpoint ===
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 API Alumni & Pekerjaan (MongoDB version) berjalan!")
	})

	return app
}
