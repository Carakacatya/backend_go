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
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment")
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("❌ DB_DSN tidak ditemukan di .env")
	}

	// Koneksi ke database
	db := database.ConnectDB(dsn)

	// Buat Fiber app
	app := config.NewApp(db)

	// Register routes (praktikum 3 & 4)
	route.AlumniRoute(app, db)
	route.PekerjaanRoute(app, db)
	route.AlumniStatusRoute(app)

	// Register route auth (modul 5)
	route.AuthRoute(app, db)

	// Jalankan server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server running at http://127.0.0.1:%s", port)
	log.Fatal(app.Listen(":" + port))
}
