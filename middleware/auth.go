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
		token := ""
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

		// Konversi string UserID dari JWT ke ObjectID (jika valid)
		userID := claims.UserID
		if _, err := primitive.ObjectIDFromHex(userID); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "ID user tidak valid",
			})
		}

		// Simpan user ke context sebagai map[string]interface{}
		c.Locals("user", map[string]interface{}{
			"id":       userID,
			"username": claims.Username,
			"role":     claims.Role,
		})
		c.Locals("user_id", userID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// AdminOnly middleware untuk admin saja
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

// UserOnly middleware untuk user biasa saja
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
