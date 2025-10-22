package config

import (
	"database/sql"
	"encoding/json"

	"praktikum3/middleware"
	"praktikum3/route"
	"praktikum3/app/repository"
	"praktikum3/app/service"

	"github.com/gofiber/fiber/v2"
)

func NewApp(db *sql.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ")
		},
		JSONDecoder: json.Unmarshal,
	})

	app.Use(middleware.LoggerMiddleware)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 API Alumni & Pekerjaan berjalan. Coba /praktikum3/alumni atau /praktikum3/pekerjaan")
	})

	api := app.Group("/praktikum3")
	route.AlumniRoute(api, db)
	route.PekerjaanRoute(api, db)

	alumniRepo := repository.NewAlumniRepository(db)
	pekerjaanRepo := repository.NewPekerjaanRepository(db)

	alumniService := service.NewAlumniService(alumniRepo)
	pekerjaanService := service.NewPekerjaanService(pekerjaanRepo)

	api = app.Group("/api")
	route.AuthRoute(api, db)

	alumni := api.Group("/alumni", middleware.AuthRequired())
	alumni.Get("/", alumniService.GetAll)
	alumni.Get("/:id", alumniService.GetByID)
	alumni.Post("/", middleware.AdminOnly(), alumniService.Create)
	alumni.Put("/:id", middleware.AdminOnly(), alumniService.Update)
	alumni.Delete("/:id", middleware.AdminOnly(), alumniService.SoftDelete)

	pekerjaan := api.Group("/pekerjaan", middleware.AuthRequired())
	// pekerjaan.Get("/", pekerjaanService.GetAll)
	pekerjaan.Get("/:id", pekerjaanService.GetByID)
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), pekerjaanService.GetByAlumniID)
	pekerjaan.Post("/", middleware.AdminOnly(), pekerjaanService.Create)
	pekerjaan.Put("/:id", middleware.AdminOnly(), pekerjaanService.Update)
	pekerjaan.Delete("/:id", middleware.AdminOnly(), pekerjaanService.SoftDelete)

	user := app.Group("/user")
	user.Delete("/pekerjaan/:id", pekerjaanService.SoftDelete)

	// user.Delete("/:id", service.DeleteUserHandler(db))


	// user := api.Group("/user", middleware.AuthRequired())
	// user.Delete("/:id", func(c *fiber.Ctx) error {
	// 	role := c.Locals("role").(string)
	// 	id, err := c.ParamsInt("id")
	// 	if err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"success": false,
	// 			"error":   "invalid user id",
	// 		})
	// 	}

	// 	if err := service.SoftDeleteUser(db, uint(id), role); err != nil {
	// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 			"success": false,
	// 			"error":   err.Error(),
	// 		})
	// 	}

	// 	return c.JSON(fiber.Map{
	// 		"success": true,
	// 		"message": "User berhasil dihapus (soft delete)",
	// 	})
	// })

	return app
}
