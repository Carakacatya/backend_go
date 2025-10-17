package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connected successfully")
	
	DB = db

	return db
}
