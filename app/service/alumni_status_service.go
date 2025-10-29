package service

import (
	"praktikum3/app/repository"

	"github.com/gofiber/fiber/v2"
)

type AlumniStatusService struct {
	statusRepo repository.AlumniStatusRepository
}

// Constructor
func NewAlumniStatusService(repo repository.AlumniStatusRepository) *AlumniStatusService {
	return &AlumniStatusService{statusRepo: repo}
}

// ✅ Handler untuk mendapatkan laporan berdasarkan status pekerjaan
func (s *AlumniStatusService) GetAlumniByStatus(c *fiber.Ctx) error {
	status := c.Query("status", "aktif") // default: aktif

	data, count, err := s.statusRepo.GetAlumniByStatus(status)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"status":   status,
		"jumlah_bekerja_lebih_dari_satu_tahun": count,
		"data":     data,
	})
}
