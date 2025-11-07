package service

import (
	"mime/multipart"
	"os"
	"path/filepath"

	"praktikum3/app/model"
	"praktikum3/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileService struct {
	repo       repository.FileRepository
	uploadBase string
}

func NewFileService(repo repository.FileRepository, uploadBase string) *FileService {
	return &FileService{repo: repo, uploadBase: uploadBase}
}

// helper: simpan file ke disk
func (s *FileService) saveFileToDisk(c *fiber.Ctx, fh *multipart.FileHeader, folder string) (string, error) {
	ext := filepath.Ext(fh.Filename)
	newName := uuid.New().String() + ext
	dir := filepath.Join(s.uploadBase, folder)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	path := filepath.Join(dir, newName)

	if err := c.SaveFile(fh, path); err != nil {
		return "", err
	}
	return path, nil
}

// === Upload foto (max 1MB, jpg/png)
func (s *FileService) UploadPhoto(c *fiber.Ctx) error {
	return s.uploadHandler(c, "photo", []string{"image/jpeg", "image/png", "image/jpg"}, 1*1024*1024)
}

// === Upload sertifikat (max 2MB, pdf)
func (s *FileService) UploadCertificate(c *fiber.Ctx) error {
	return s.uploadHandler(c, "certificate", []string{"application/pdf"}, 2*1024*1024)
}

func (s *FileService) uploadHandler(c *fiber.Ctx, category string, allowed []string, maxBytes int64) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "file tidak ditemukan",
			"error":   err.Error(),
		})
	}

	if fh.Size > maxBytes {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ukuran file melebihi batas",
		})
	}

	contentType := fh.Header.Get("Content-Type")
	ok := false
	for _, t := range allowed {
		if t == contentType {
			ok = true
			break
		}
	}
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "format file tidak diizinkan",
		})
	}

	// === Ambil user dari context ===
	var userID primitive.ObjectID
	var role string

	userAny := c.Locals("user")
	if userAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid token data",
		})
	}

	switch u := userAny.(type) {
	case *model.User:
		userID = u.ID
		role = u.Role
	case map[string]interface{}:
		if idStr, ok := u["id"].(string); ok {
			uid, _ := primitive.ObjectIDFromHex(idStr)
			userID = uid
		}
		if r, ok := u["role"].(string); ok {
			role = r
		}
	default:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User context tidak valid",
		})
	}

	// === Admin bisa upload untuk user lain ===
	var alumniObj *primitive.ObjectID
	alumniIDForm := c.FormValue("alumni_id", "")
	if alumniIDForm != "" {
		oid, err := primitive.ObjectIDFromHex(alumniIDForm)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "alumni_id tidak valid",
			})
		}
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "hanya admin boleh upload untuk user lain",
			})
		}
		alumniObj = &oid
	} else {
		alumniObj = &userID
	}

	folder := map[string]string{
		"photo":       "photos",
		"certificate": "certificates",
	}[category]

	path, err := s.saveFileToDisk(c, fh, folder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "gagal menyimpan file",
			"error":   err.Error(),
		})
	}

	f := &model.File{
		AlumniID:     alumniObj,
		FileName:     filepath.Base(path),
		OriginalName: fh.Filename,
		FilePath:     path,
		FileSize:     fh.Size,
		FileType:     contentType,
		Category:     category,
		UploadedBy:   userID,
	}

	id, err := s.repo.Create(f)
	if err != nil {
		os.Remove(path)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "gagal menyimpan metadata",
			"error":   err.Error(),
		})
	}
	f.ID = id

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "file uploaded",
		"data":    f,
	})
}

// === GET /api/files
func (s *FileService) GetAll(c *fiber.Ctx) error {
	list, err := s.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{"success": true, "data": list})
}

// === GET /api/files/:id
func (s *FileService) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "id tidak valid",
		})
	}
	f, err := s.repo.FindByID(oid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "file tidak ditemukan",
		})
	}
	return c.JSON(fiber.Map{"success": true, "data": f})
}

// === DELETE /api/files/:id
func (s *FileService) DeleteByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "id tidak valid",
		})
	}

	f, err := s.repo.FindByID(oid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "file tidak ditemukan",
		})
	}

	userAny := c.Locals("user")
	if userAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid token data",
		})
	}

	var userID primitive.ObjectID
	var role string

	switch u := userAny.(type) {
	case *model.User:
		userID = u.ID
		role = u.Role
	case map[string]interface{}:
		if idStr, ok := u["id"].(string); ok {
			userID, _ = primitive.ObjectIDFromHex(idStr)
		}
		if r, ok := u["role"].(string); ok {
			role = r
		}
	}

	if role != "admin" && f.UploadedBy != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "tidak punya akses menghapus file ini",
		})
	}

	os.Remove(f.FilePath)
	if err := s.repo.DeleteByID(oid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "gagal hapus metadata",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true, "message": "file dihapus"})
}
