package route

import (
	"database/sql"
	"praktikum3/app/repository"
	"praktikum3/app/service"
	"praktikum3/middleware"

	"github.com/gofiber/fiber/v2"
)

func PekerjaanRoute(r fiber.Router, db *sql.DB) {
	repo := repository.NewPekerjaanRepository(db)
	p := service.NewPekerjaanService(repo)

	g := r.Group("/pekerjaan", middleware.AuthRequired())

	// g.Get("/", p.GetAll)
	g.Get("/trash", middleware.AdminOnly(), p.GetTrashed)
	g.Get("/:id", p.GetByID)
	g.Get("/alumni/:alumni_id", p.GetByAlumniID)
	g.Post("/", middleware.AdminOnly(), p.Create)
	g.Put("/:id", middleware.AdminOnly(), p.Update)
	g.Delete("/:id", p.SoftDelete)
	g.Put("/restore/:id", p.Restore)
	g.Delete("/hard/:id", middleware.AdminOnly(), p.HardDelete)
}
