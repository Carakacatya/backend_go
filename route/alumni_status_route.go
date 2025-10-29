package route

import (
	"praktikum3/app/repository"
	"praktikum3/app/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// AlumniStatusRoute mendaftarkan semua endpoint untuk data alumni berdasarkan status pekerjaan
func AlumniStatusRoute(app *fiber.App, db *mongo.Database) {
	// Inisialisasi repository dan service
	statusRepo := repository.NewAlumniStatusRepository(db)
	statusService := service.NewAlumniStatusService(statusRepo)

	// Buat group route untuk alumni status
	r := app.Group("/alumni-status")

	// GET /alumni-status?status=aktif
	r.Get("/", statusService.GetAlumniByStatus)
}
