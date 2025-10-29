package route

import (
	"praktikum3/app/repository"
	"praktikum3/app/service"
	"praktikum3/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func PekerjaanRoute(r fiber.Router, db *mongo.Database) {
	repo := repository.NewPekerjaanRepository(db)
	p := service.NewPekerjaanService(repo)

	g := r.Group("/pekerjaan", middleware.AuthRequired())

	// Admin only
	g.Post("/", middleware.AdminOnly(), p.Create)
	g.Put("/:id", middleware.AdminOnly(), p.Update)
	g.Delete("/hard/:id", middleware.AdminOnly(), p.HardDelete)
	g.Get("/trash", middleware.AdminOnly(), p.GetTrashed)

	// Semua user
	g.Get("/", p.GetAll)
	g.Get("/:id", p.GetByID)
	g.Get("/alumni/:alumni_id", p.GetByAlumniID)
	g.Delete("/:id", p.SoftDelete)
	g.Put("/restore/:id", middleware.AdminOnly(), p.Restore)
}
