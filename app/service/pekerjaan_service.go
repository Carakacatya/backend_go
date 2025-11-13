package service

import (
	"time"

	"praktikum3/app/model"
	"praktikum3/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PekerjaanService struct {
	repo repository.PekerjaanRepository
}

func NewPekerjaanService(repo repository.PekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{repo: repo}
}

// ================== GET ALL ==================
// GetAll godoc
// @Summary Get semua pekerjaan
// @Description Mengambil semua data pekerjaan alumni tanpa parameter
// @Tags Pekerjaan
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /pekerjaan/ [get]
func (s *PekerjaanService) GetAll(c *fiber.Ctx) error {
	search := c.Query("search", "")
	sortBy := c.Query("sortBy", "created_at")
	order := c.Query("order", "DESC")
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	offset := (page - 1) * limit

	data, err := s.repo.GetAllWithQuery(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	count, err := s.repo.Count(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"meta": fiber.Map{
			"total": count,
			"page":  page,
			"limit": limit,
		},
	})
}

// ================== GET BY ID ==================
// @Summary Get pekerjaan by ID
// @Description Mendapatkan detail pekerjaan berdasarkan ID
// @Tags Pekerjaan
// @Security BearerAuth
// @Param id path string true "ID Pekerjaan"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /pekerjaan/{id} [get]
func (s *PekerjaanService) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	data, err := s.repo.GetByID(objectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	if data == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"success": true, "data": data})
}

// ================== GET BY ALUMNI ID ==================
// @Summary Get pekerjaan berdasarkan alumni ID
// @Description Mendapatkan semua pekerjaan milik alumni tertentu
// @Tags Pekerjaan
// @Security BearerAuth
// @Param alumni_id path string true "ID Alumni"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]interface{}
// @Router /pekerjaan/alumni/{alumni_id} [get]
func (s *PekerjaanService) GetByAlumniID(c *fiber.Ctx) error {
	alumniIDParam := c.Params("alumni_id")
	alumniID, err := primitive.ObjectIDFromHex(alumniIDParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "alumni_id tidak valid"})
	}

	userData := c.Locals("user")
	claimsMap, ok := userData.(map[string]interface{})
	if !ok {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid token data"})
	}

	role, _ := claimsMap["role"].(string)
	includeDeleted := role == "admin"

	data, err := s.repo.GetByAlumniID(alumniID, includeDeleted)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": data})
}

// ================== CREATE ==================
// Create godoc
// @Summary Tambah pekerjaan alumni
// @Description Tambahkan data pekerjaan untuk alumni tertentu (wajib kirim alumni_id)
// @Tags Pekerjaan
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.CreatePekerjaanReq true "Data pekerjaan untuk alumni tertentu"
// @Success 200 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]interface{}
// @Router /pekerjaan/ [post]
func (s *PekerjaanService) Create(c *fiber.Ctx) error {
	var in model.CreatePekerjaanReq
	if err := c.BodyParser(&in); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	start, err := time.Parse("2006-01-02", in.TanggalMulaiKerja)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Tanggal mulai tidak valid"})
	}

	var end *time.Time
	if in.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", in.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "Tanggal selesai tidak valid"})
		}
		end = &t
	}

	id, err := s.repo.Create(in, &start, end)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dibuat", "id": id.Hex()})
}

// ================== UPDATE ==================
// @Summary Update pekerjaan
// @Description Mengupdate data pekerjaan berdasarkan ID
// @Tags Pekerjaan
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID Pekerjaan"
// @Param pekerjaan body model.UpdatePekerjaanReq true "Data pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404,500 {object} map[string]interface{}
// @Router /pekerjaan/{id} [put]
func (s *PekerjaanService) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID pekerjaan tidak valid"})
	}

	var in model.UpdatePekerjaanReq
	if err := c.BodyParser(&in); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	start, err := time.Parse("2006-01-02", in.TanggalMulaiKerja)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Tanggal mulai tidak valid"})
	}

	var end *time.Time
	if in.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", in.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "Tanggal selesai tidak valid"})
		}
		end = &t
	}

	if err := s.repo.Update(objectID, in, &start, end); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil diperbarui"})
}

// ================== SOFT DELETE ==================
// @Summary Soft delete pekerjaan
// @Description Menghapus pekerjaan tanpa menghapus permanen
// @Tags Pekerjaan
// @Security BearerAuth
// @Param id path string true "ID Pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]interface{}
// @Router /pekerjaan/{id} [delete]
func (s *PekerjaanService) SoftDelete(c *fiber.Ctx) error {
	userData := c.Locals("user")
	claimsMap, ok := userData.(map[string]interface{})
	if !ok {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid token data"})
	}

	role, _ := claimsMap["role"].(string)
	userIDStr, _ := claimsMap["id"].(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "User ID tidak valid"})
	}

	idStr := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID pekerjaan tidak valid"})
	}

	if role == "admin" {
		err = s.repo.SoftDeleteByAdmin(objectID)
	} else {
		err = s.repo.SoftDeleteByUser(objectID, userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan dihapus (soft delete)"})
}

// ================== RESTORE ==================
// @Summary Restore pekerjaan
// @Description Mengembalikan pekerjaan yang dihapus (soft delete)
// @Tags Pekerjaan
// @Security BearerAuth
// @Param id path string true "ID Pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]interface{}
// @Router /pekerjaan/restore/{id} [put]
func (s *PekerjaanService) Restore(c *fiber.Ctx) error {
	idStr := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	if err := s.repo.RestoreByID(objectID); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data berhasil direstore"})
}

// ================== HARD DELETE ==================
// @Summary Hard delete pekerjaan
// @Description Menghapus data pekerjaan secara permanen
// @Tags Pekerjaan
// @Security BearerAuth
// @Param id path string true "ID Pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]interface{}
// @Router /pekerjaan/hard/{id} [delete]
func (s *PekerjaanService) HardDelete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	if err := s.repo.HardDeleteByID(objectID); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data dihapus permanen"})
}

// ================== GET TRASH ==================
// @Summary Get pekerjaan yang dihapus (trash)
// @Description Mengambil semua pekerjaan yang dihapus (soft delete)
// @Tags Pekerjaan
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /pekerjaan/trash [get]
func (s *PekerjaanService) GetTrashed(c *fiber.Ctx) error {
	data, err := s.repo.GetAllTrash()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": data})
}
