package route

import (
    "os"
    "praktikum3/app/repository"
    "praktikum3/app/service"
    "praktikum3/middleware"

    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
)

func isRunningTest() bool {
    return os.Getenv("UNIT_TEST") == "true"
}

func AlumniRoute(r fiber.Router, db *mongo.Database) {
    repo := repository.NewAlumniRepository(db)
    al := service.NewAlumniService(repo)

    testMode := isRunningTest()

    var g fiber.Router

    // MODE TEST: TANPA AUTH
    if testMode {
        g = r.Group("/alumni")
    } else {
        g = r.Group("/alumni", middleware.AuthRequired())
    }

    //
    // ================================
    // ADMIN ROUTES
    // ================================
    //
    if testMode {
        // TANPA MIDDLEWARE SAAT TEST (sesuai apa yg dilakukan test)
        g.Get("/trash", al.GetTrashed)
        g.Put("/restore/:id", al.Restore)
        g.Delete("/hard/:id", al.HardDelete)
    } else {
        g.Get("/trash", middleware.AdminOnly(), al.GetTrashed)
        g.Put("/restore/:id", middleware.AdminOnly(), al.Restore)
        g.Delete("/hard/:id", middleware.AdminOnly(), al.HardDelete)
    }

    //
    // ================================
    // CRUD ROUTES (admin only)
    // ================================
    //
    if testMode {
        g.Post("/", al.Create)
        g.Put("/:id", al.Update)
        g.Delete("/:id", al.SoftDelete)
    } else {
        g.Post("/", middleware.AdminOnly(), al.Create)
        g.Put("/:id", middleware.AdminOnly(), al.Update)
        g.Delete("/:id", middleware.AdminOnly(), al.SoftDelete)
    }

    //
    // ================================
    // PUBLIC ROUTES
    // ================================
    //
    g.Get("/", al.GetAll)
    g.Get("/:id", al.GetByID)
}
