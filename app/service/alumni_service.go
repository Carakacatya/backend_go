package service

import (
	"praktikum3/app/model"
	"praktikum3/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlumniService struct {
	alumniRepo repository.AlumniRepository
}

func NewAlumniService(repo repository.AlumniRepository) *AlumniService {
	return &AlumniService{alumniRepo: repo}
}

// 🔹 Get All Alumni
func (s *AlumniService) GetAll(c *fiber.Ctx) error {
	data, err := s.alumniRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

// 🔹 Get Alumni by ID
func (s *AlumniService) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	data, err := s.alumniRepo.GetByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if data == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

// 🔹 Create Alumni
func (s *AlumniService) Create(c *fiber.Ctx) error {
	var alumni model.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if alumni.Nama == "" || alumni.Email == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Nama dan Email wajib diisi"})
	}

	err := s.alumniRepo.Create(&alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Alumni berhasil ditambahkan"})
}

// 🔹 Update Alumni
func (s *AlumniService) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	var alumni model.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	existing, err := s.alumniRepo.GetByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if existing == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Alumni tidak ditemukan"})
	}

	err = s.alumniRepo.Update(id, &alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data berhasil diperbarui"})
}

// 🔹 Soft Delete Alumni
func (s *AlumniService) SoftDelete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	err = s.alumniRepo.SoftDelete(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data berhasil dihapus (soft delete)"})
}

// 🔹 Get Trashed Alumni
func (s *AlumniService) GetTrashed(c *fiber.Ctx) error {
	data, err := s.alumniRepo.GetTrashed()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

// 🔹 Restore Alumni
func (s *AlumniService) Restore(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	existing, err := s.alumniRepo.GetTrashedByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if existing == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data tidak ditemukan di trash"})
	}

	err = s.alumniRepo.Restore(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data berhasil direstore"})
}

// 🔹 Hard Delete Alumni
func (s *AlumniService) HardDelete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	existing, err := s.alumniRepo.GetTrashedByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if existing == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data tidak ditemukan di trash"})
	}

	err = s.alumniRepo.ForceDelete(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Data dihapus permanen"})
}
