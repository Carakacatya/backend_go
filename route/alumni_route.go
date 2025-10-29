package route

import (
	"praktikum3/app/repository"
	"praktikum3/app/service"
	"praktikum3/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// AlumniRoute untuk endpoint CRUD Alumni berbasis MongoDB
func AlumniRoute(r fiber.Router, db *mongo.Database) {
	repo := repository.NewAlumniRepository(db)
	al := service.NewAlumniService(repo)

	g := r.Group("/alumni", middleware.AuthRequired())

	// Admin-only features
	g.Get("/trash", middleware.AdminOnly(), al.GetTrashed)
	g.Put("/restore/:id", middleware.AdminOnly(), al.Restore)
	g.Delete("/hard/:id", middleware.AdminOnly(), al.HardDelete)

	// General CRUD
	g.Get("/", al.GetAll)
	g.Get("/:id", al.GetByID)
	g.Post("/", middleware.AdminOnly(), al.Create)
	g.Put("/:id", middleware.AdminOnly(), al.Update)
	g.Delete("/:id", middleware.AdminOnly(), al.SoftDelete)
}
