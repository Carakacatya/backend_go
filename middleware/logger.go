package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware menampilkan log setiap request
func LoggerMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	duration := time.Since(start)

	statusColor := ""
	resetColor := "\033[0m"

	switch {
	case c.Response().StatusCode() >= 500:
		statusColor = "\033[31m" // merah
	case c.Response().StatusCode() >= 400:
		statusColor = "\033[33m" // kuning
	case c.Response().StatusCode() >= 300:
		statusColor = "\033[36m" // biru muda
	default:
		statusColor = "\033[32m" // hijau
	}

	fmt.Printf("%s[%d]%s %-7s %s (%v)\n",
		statusColor,
		c.Response().StatusCode(),
		resetColor,
		c.Method(),
		c.Path(),
		duration,
	)

	return err
}
