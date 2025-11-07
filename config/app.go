package config

import (
	"encoding/json"
	"log"
	"os"
	"praktikum3/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewApp membuat instance Fiber dengan konfigurasi global
func NewApp(db *mongo.Database) *fiber.App {
	app := fiber.New(fiber.Config{
		// ✅ Batas maksimal ukuran body agar upload file besar bisa
		BodyLimit: 10 * 1024 * 1024, // 10MB (untuk upload foto/pdf)
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ") // biar JSON rapi
		},
		JSONDecoder: json.Unmarshal,
	})

	// === 1️⃣ Pastikan folder upload sudah ada ===
	ensureUploadDirs()

	// === 2️⃣ Middleware global ===
	app.Use(middleware.LoggerMiddleware)

	// === 3️⃣ Static file route (akses langsung ke file upload) ===
	// contoh akses: http://localhost:3000/uploads/photos/nama.jpg
	app.Static("/uploads", "./uploads")

	// === 4️⃣ Root endpoint ===
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 API Alumni & Pekerjaan (MongoDB + Upload File) berjalan dengan baik!")
	})

	return app
}

// ensureUploadDirs memastikan folder upload tersedia saat server pertama dijalankan
func ensureUploadDirs() {
	baseDirs := []string{
		"./uploads/photos",
		"./uploads/certificates",
	}

	for _, dir := range baseDirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("❌ Gagal membuat folder upload: %v", err)
		}
	}
}
