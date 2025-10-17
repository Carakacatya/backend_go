package service

import (
	"praktikum3/app/repository"

	"github.com/gofiber/fiber/v2"
)

func GetAlumniByStatusService(c *fiber.Ctx) error {
    status := c.Query("status", "aktif") // default aktif

    data, count, err := repository.GetAlumniByStatus(status)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Gagal mengambil data",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "status": status,
        "jumlah_bekerja_lebih_dari_satu_tahun": count,
        "data": data,
    })
}
