package service

import (
	"praktikum3/app/model"
	"praktikum3/app/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanService struct {
	repo *repository.PekerjaanRepository
}

func NewPekerjaanService(repo *repository.PekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{repo: repo}
}

func (s *PekerjaanService) GetAll(c *fiber.Ctx) error {
	search := c.Query("search", "")
	sortBy := c.Query("sortBy", "created_at")
	order := c.Query("order", "DESC")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	data, err := s.repo.GetAllWithQuery(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	total, _ := s.repo.Count(search)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (s *PekerjaanService) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	data, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) GetByAlumniID(c *fiber.Ctx) error {
	alumniID, _ := strconv.Atoi(c.Params("alumni_id"))
	data, err := s.repo.GetByAlumniID(alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) Create(c *fiber.Ctx) error {
	var in model.CreatePekerjaanReq
	if err := c.BodyParser(&in); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	start, _ := time.Parse("2006-01-02", in.TanggalMulaiKerja)
	var end *time.Time
	if in.TanggalSelesaiKerja != "" {
		t, _ := time.Parse("2006-01-02", in.TanggalSelesaiKerja)
		end = &t
	}
	id, err := s.repo.Create(in, &start, end)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dibuat", "id": id})
}

func (s *PekerjaanService) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var in model.UpdatePekerjaanReq
	if err := c.BodyParser(&in); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	start, _ := time.Parse("2006-01-02", in.TanggalMulaiKerja)
	var end *time.Time
	if in.TanggalSelesaiKerja != "" {
		t, _ := time.Parse("2006-01-02", in.TanggalSelesaiKerja)
		end = &t
	}
	err := s.repo.Update(id, in, &start, end)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil diupdate"})
}

func (s *PekerjaanService) SoftDelete(c *fiber.Ctx) error {
	userData := c.Locals("user")
	claimsMap, ok := userData.(map[string]interface{})
	if !ok {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid token data"})
	}

	role, _ := claimsMap["role"].(string)

	var userID int
	if v, ok := claimsMap["id"].(float64); ok {
		userID = int(v)
	} else if v, ok := claimsMap["id"].(int); ok {
		userID = v
	} else {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid user ID type in token"})
	}

	id, _ := strconv.Atoi(c.Params("id"))

	if role == "admin" {
		err := s.repo.SoftDeleteByID(id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
		}
	} else {
		err := s.repo.SoftDeleteByIDAndUser(id, userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
		}
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan dihapus (soft delete)"})
}

func (s *PekerjaanService) GetTrashed(c *fiber.Ctx) error {
	userData := c.Locals("user")
	claimsMap, ok := userData.(map[string]interface{})
	if !ok {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid token data"})
	}

	role, _ := claimsMap["role"].(string)
	var userID int
	if v, ok := claimsMap["id"].(float64); ok {
		userID = int(v)
	} else if v, ok := claimsMap["id"].(int); ok {
		userID = v
	} else {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid user ID type in token"})
	}

	var data []model.PekerjaanAlumni
	var err error

	if role == "admin" {
		data, err = s.repo.GetAllTrash()
	} else {
		data, err = s.repo.GetUserTrash(userID)
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

func (s *PekerjaanService) Restore(c *fiber.Ctx) error {
	userData := c.Locals("user")
	claimsMap, ok := userData.(map[string]interface{})
	if !ok {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid token data"})
	}

	role, _ := claimsMap["role"].(string)
	var userID int
	if v, ok := claimsMap["id"].(float64); ok {
		userID = int(v)
	} else if v, ok := claimsMap["id"].(int); ok {
		userID = v
	} else {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid user ID type in token"})
	}

	id, _ := strconv.Atoi(c.Params("id"))
	var err error

	if role == "admin" {
		err = s.repo.RestoreByID(id)
	} else {
		err = s.repo.RestoreByIDAndUser(id, userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data berhasil direstore"})
}

func (s *PekerjaanService) HardDelete(c *fiber.Ctx) error {
	userData := c.Locals("user")
	claimsMap, ok := userData.(map[string]interface{})
	if !ok {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid token data"})
	}

	role, _ := claimsMap["role"].(string)
	var userID int
	if v, ok := claimsMap["id"].(float64); ok {
		userID = int(v)
	} else if v, ok := claimsMap["id"].(int); ok {
		userID = v
	} else {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Invalid user ID type in token"})
	}

	id, _ := strconv.Atoi(c.Params("id"))
	var err error

	if role == "admin" {
		err = s.repo.HardDeleteByID(id)
	} else {
		err = s.repo.HardDeleteByUser(id, userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data dihapus permanen"})
}
