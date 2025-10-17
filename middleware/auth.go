package middleware

import (
	"praktikum3/app/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Authorization token diperlukan",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Format token salah, gunakan 'Bearer <token>'",
			})
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid atau sudah expired",
			})
		}

		userMap := map[string]any{
			"id":       claims.UserID,
			"username": claims.Username,
			"role":     claims.Role,
		}
		c.Locals("user", userMap)
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Hanya admin yang boleh mengakses endpoint ini",
			})
		}
		return c.Next()
	}
}

func UserOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role != "user" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Hanya user yang boleh mengakses endpoint ini",
			})
		}
		return c.Next()
	}
}
