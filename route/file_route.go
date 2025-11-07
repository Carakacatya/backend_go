package route

import (
	"praktikum3/app/repository"
	"praktikum3/app/service"
	"praktikum3/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func FileRoute(app fiber.Router, db *mongo.Database, uploadBase string) {
	repo := repository.NewFileRepository(db)
	svc := service.NewFileService(repo, uploadBase)

	// Semua endpoint dalam /api/files wajib login
	api := app.Group("/api/files", middleware.AuthRequired())

	// === Upload ===
	// Admin dan user bisa upload, tapi validasi siapa boleh upload untuk siapa ada di service
	api.Post("/photo", svc.UploadPhoto)           // ✅ upload foto (jpg/png ≤1MB)
	api.Post("/certificate", svc.UploadCertificate) // ✅ upload sertifikat (pdf ≤2MB)

	// === Read ===
	api.Get("/", middleware.AdminOnly(), svc.GetAll) // ✅ hanya admin bisa lihat semua file
	api.Get("/:id", svc.GetByID)                     // ✅ semua user bisa lihat file by id

	// === Delete ===
	api.Delete("/:id", svc.DeleteByID) // ✅ admin bisa hapus semua, user hanya miliknya sendiri
}
