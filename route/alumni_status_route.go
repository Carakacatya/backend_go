package route

import (
	"praktikum3/app/service"

	"github.com/gofiber/fiber/v2"
)

func AlumniStatusRoute(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/laporan/alumni-by-status", service.GetAlumniByStatusService)
}
