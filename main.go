package main

import (
	"log"
	"os"

	"praktikum3/config"
	"praktikum3/database"
	"praktikum3/route"

	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger" // swagger middleware fiber
	_ "praktikum3/docs"                            // import docs swagger
)

//
// ========== SWAGGER INFO ==========
//

// @title Alumni API Documentation
// @version 1.0
// @description API untuk mengelola data alumni dengan MongoDB dan Clean Architecture
// @host localhost:3000
// @BasePath /api/v1
// @schemes http

// ✅ JWT Security Definition untuk Swagger
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// === 1️⃣ Load environment variables ===
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	// === 2️⃣ Connect to MongoDB ===
	mongoDB := database.ConnectMongo()
	if mongoDB == nil {
		log.Fatal("❌ Failed to connect to MongoDB")
	}

	// === 3️⃣ Initialize Fiber App ===
	app := config.NewApp(mongoDB)

	// ✅ Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// === 4️⃣ Static Files (upload folder) ===
	app.Static("/uploads", "./uploads")

	// === 5️⃣ ROUTES ===
	api := app.Group("/api/v1")

	route.AuthRoute(api, mongoDB)
	route.AlumniRoute(api, mongoDB)
	route.PekerjaanRoute(api, mongoDB)
	route.AlumniStatusRoute(app, mongoDB) // ini tidak di bawah /api/v1
	route.FileRoute(api, mongoDB, "./uploads")

	// === 6️⃣ PORT ===
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// === 7️⃣ RUN SERVER ===
	log.Printf("🚀 Server running at http://127.0.0.1:%s", port)
	log.Println("📄 Swagger running at http://localhost:" + port + "/swagger/index.html")

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
