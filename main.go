package main

import (
	"log"
	"os"

	"praktikum3/config"
	"praktikum3/database"
	"praktikum3/route"

	"github.com/joho/godotenv"
)

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

	// === 4️⃣ Serve static files (untuk akses file upload) ===
	// Contoh akses: http://localhost:3000/uploads/photos/nama.jpg
	app.Static("/uploads", "./uploads")

	// === 5️⃣ Register semua routes utama ===
	route.AuthRoute(app, mongoDB)              // login / register
	route.AlumniRoute(app, mongoDB)            // data alumni
	route.PekerjaanRoute(app, mongoDB)         // data pekerjaan alumni
	route.AlumniStatusRoute(app, mongoDB)      // status alumni
	route.FileRoute(app, mongoDB, "./uploads") // ✅ upload foto & sertifikat

	// === 6️⃣ Get port dari environment (.env) ===
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// === 7️⃣ Jalankan server ===
	log.Printf("🚀 Server running at http://127.0.0.1:%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
