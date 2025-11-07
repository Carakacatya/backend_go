package middleware

import (
	"praktikum3/app/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthRequired middleware untuk endpoint yang wajib login
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Authorization token diperlukan",
			})
		}

		// Ambil token tanpa kata "Bearer "
		var token string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Format token salah, gunakan 'Bearer <token>'",
			})
		}

		// Validasi token JWT
		claims, err := utils.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid atau sudah expired",
			})
		}

		// Validasi ObjectID
		if _, err := primitive.ObjectIDFromHex(claims.UserID); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "ID user tidak valid",
			})
		}

		// Simpan data user ke context
		c.Locals("user", map[string]interface{}{
			"id":       claims.UserID,
			"username": claims.Username,
			"role":     claims.Role,
		})
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// AdminOnly middleware untuk admin saja
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Hanya admin yang boleh mengakses endpoint ini",
			})
		}
		return c.Next()
	}
}

// UserOnly middleware untuk user biasa saja
func UserOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "user" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Hanya user yang boleh mengakses endpoint ini",
			})
		}
		return c.Next()
	}
}
