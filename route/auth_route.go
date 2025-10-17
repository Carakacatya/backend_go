package route

import (
	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app fiber.Router) {
	app.Post("/login", func(c *fiber.Ctx) error {
		return c.SendString("login not implemented")
	})
}
