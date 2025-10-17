package config

import (
	"log"
	"os"
)

// InitLogger membuat logger baru yang bisa digunakan di seluruh aplikasi
func InitLogger() *log.Logger {
	// Output logger diarahkan ke stdout (console)
	return log.New(os.Stdout, "[PRAKTIKUM3] ", log.Ldate|log.Ltime|log.Lshortfile)
}
