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
	// === Load environment variables ===
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	// === Connect to MongoDB ===
	mongoDB := database.ConnectMongo()
	if mongoDB == nil {
		log.Fatal("❌ Failed to connect to MongoDB")
	}

	// === Initialize Fiber App ===
	app := config.NewApp(mongoDB)

	// === Register Routes ===
	route.AuthRoute(app, mongoDB)
	route.AlumniRoute(app, mongoDB)
	route.PekerjaanRoute(app, mongoDB)
	route.AlumniStatusRoute(app, mongoDB) // ✅ pastikan route ini juga menerima db jika dibutuhkan

	// === Get port from environment ===
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// === Start the server ===
	log.Printf("🚀 Server running at http://127.0.0.1:%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
