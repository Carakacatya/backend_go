package route

import (
	"database/sql"
	"praktikum3/app/repository"
	"praktikum3/app/service"
	"praktikum3/middleware"

	"github.com/gofiber/fiber/v2"
)

func AlumniRoute(r fiber.Router, db *sql.DB) {
	repo := repository.NewAlumniRepository(db)
	al := service.NewAlumniService(repo)

	g := r.Group("/alumni", middleware.AuthRequired())

	g.Get("/trash", middleware.AdminOnly(), al.GetTrashed)
	g.Put("/restore/:id", al.Restore)
	g.Delete("/hard/:id", middleware.AdminOnly(), al.HardDelete)

	g.Get("/", al.GetAll)
	g.Get("/:id", al.GetByID)
	g.Post("/", middleware.AdminOnly(), al.Create)
	g.Put("/:id", middleware.AdminOnly(), al.Update)
	g.Delete("/:id", middleware.AdminOnly(), al.SoftDelete)
}
