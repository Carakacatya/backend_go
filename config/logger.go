package config

import (
	"log"
	"os"
)

// InitLogger membuat instance logger global untuk aplikasi.
// Logger ini menulis ke stdout dengan format tanggal, waktu, dan nama file baris kode.
func InitLogger() *log.Logger {
	return log.New(os.Stdout, "📘 [PRAKTIKUM3] ", log.Ldate|log.Ltime|log.Lshortfile)
}
