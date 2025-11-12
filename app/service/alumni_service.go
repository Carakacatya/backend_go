package service

import (
	"time"

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

// GetAll godoc
// @Summary Get semua alumni
// @Description Mengambil semua data alumni aktif
// @Tags Alumni
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /alumni/ [get]
func (s *AlumniService) GetAll(c *fiber.Ctx) error {
	data, err := s.alumniRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

// GetByID godoc
// @Summary Get alumni by ID
// @Description Mendapatkan detail alumni berdasarkan ID
// @Tags Alumni
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID Alumni"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /alumni/{id} [get]
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

// Create godoc
// @Summary Tambah alumni
// @Description Menambahkan data alumni baru
// @Tags Alumni
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param alumni body model.Alumni true "Data Alumni"
// @Success 200 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]interface{}
// @Router /alumni/ [post]
func (s *AlumniService) Create(c *fiber.Ctx) error {
	var alumni model.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	if alumni.Nama == "" || alumni.Email == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Nama dan Email wajib diisi"})
	}

	// ✅ Set default value
	alumni.ID = primitive.NewObjectID()
	alumni.CreatedAt = time.Now()
	alumni.UpdatedAt = time.Now()

	err := s.alumniRepo.Create(&alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Alumni berhasil ditambahkan"})
}

// Update godoc
// @Summary Update alumni
// @Description Mengupdate data alumni berdasarkan ID
// @Tags Alumni
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID Alumni"
// @Param alumni body model.Alumni true "Data Alumni"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /alumni/{id} [put]
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

	// ✅ Update timestamp
	alumni.UpdatedAt = time.Now()

	err = s.alumniRepo.Update(id, &alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data berhasil diperbarui"})
}

// SoftDelete godoc
// @Summary Soft delete alumni
// @Description Menghapus alumni (soft delete)
// @Tags Alumni
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID Alumni"
// @Success 200 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]interface{}
// @Router /alumni/{id} [delete]
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

// GetTrashed godoc
// @Summary Get alumni yang dihapus (trash)
// @Description Mengambil alumni yang dalam status soft delete
// @Tags Alumni
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /alumni/trash [get]
func (s *AlumniService) GetTrashed(c *fiber.Ctx) error {
	data, err := s.alumniRepo.GetTrashed()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

// Restore godoc
// @Summary Restore alumni dari trash
// @Description Mengembalikan data alumni yang soft delete
// @Tags Alumni
// @Security BearerAuth
// @Param id path string true "ID Alumni"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /alumni/restore/{id} [put]
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

// HardDelete godoc
// @Summary Hard delete alumni
// @Description Menghapus data alumni secara permanen
// @Tags Alumni
// @Security BearerAuth
// @Param id path string true "ID Alumni"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /alumni/hard/{id} [delete]
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
