package service

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

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

// ====================================
// @Summary Upload foto
// @Description Upload file foto (jpg/png, max 1MB). Jika admin, wajib isi alumni_id.
// @Tags File
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File Foto"
// @Param alumni_id formData string false "ID Alumni (hanya untuk admin)"
// @Success 201 {object} map[string]interface{}
// @Failure 400,401,403,500 {object} map[string]interface{}
// @Router /api/files/photo [post]
func (s *FileService) UploadPhoto(c *fiber.Ctx) error {
	return s.uploadHandler(c, "photo", []string{"image/jpeg", "image/png", "image/jpg"}, 1*1024*1024)
}

// ====================================
// @Summary Upload sertifikat
// @Description Upload file sertifikat (PDF, max 2MB). Jika admin, wajib isi alumni_id.
// @Tags File
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File Sertifikat"
// @Param alumni_id formData string false "ID Alumni (hanya untuk admin)"
// @Success 201 {object} map[string]interface{}
// @Failure 400,401,403,500 {object} map[string]interface{}
// @Router /api/files/certificate [post]
func (s *FileService) UploadCertificate(c *fiber.Ctx) error {
	return s.uploadHandler(c, "certificate", []string{"application/pdf"}, 2*1024*1024)
}

// ====================================
// Upload Handler (helper utama)
func (s *FileService) uploadHandler(c *fiber.Ctx, category string, allowed []string, maxBytes int64) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "file tidak ditemukan",
		})
	}

	// Cek ukuran file
	if fh.Size > maxBytes {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ukuran file melebihi batas",
		})
	}

	// Validasi tipe file
	contentType := fh.Header.Get("Content-Type")
	valid := false
	for _, t := range allowed {
		if contentType == t {
			valid = true
			break
		}
	}
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "format file tidak diizinkan",
		})
	}

	// Ambil user dari JWT
	userAny := c.Locals("user")
	if userAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "User tidak ditemukan dari token",
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

	// Validasi alumni_id
	var alumniObj *primitive.ObjectID
	alumniIDForm := c.FormValue("alumni_id", "")

	if role == "admin" {
		if alumniIDForm == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "admin wajib mengisi alumni_id",
			})
		}
		oid, err := primitive.ObjectIDFromHex(alumniIDForm)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "alumni_id tidak valid",
			})
		}
		alumniObj = &oid
	} else {
		if alumniIDForm != "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "user tidak boleh menentukan alumni_id",
			})
		}
		alumniObj = &userID
	}

	// Simpan ke folder sesuai kategori
	folder := map[string]string{
		"photo":       "photos",
		"certificate": "certificates",
	}[category]

	path, err := s.saveFileToDisk(c, fh, folder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "gagal menyimpan file di server",
			"error":   err.Error(),
		})
	}

	// Buat dokumen file
	fileDoc := &model.File{
		AlumniID:     alumniObj,
		FileName:     filepath.Base(path),
		OriginalName: fh.Filename,
		FilePath:     path,
		FileSize:     fh.Size,
		FileType:     contentType,
		Category:     category,
		UploadedBy:   userID,
		UploadedAt:   time.Now(),
	}

	id, err := s.repo.Create(fileDoc)
	if err != nil {
		os.Remove(path) // rollback file kalau gagal insert DB
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "gagal menyimpan metadata file",
			"error":   err.Error(),
		})
	}
	fileDoc.ID = id

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "file berhasil diupload",
		"data":    fileDoc,
	})
}

// ====================================
// @Summary Get semua file
// @Description Mengambil semua file (khusus admin)
// @Tags File
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/files [get]
func (s *FileService) GetAll(c *fiber.Ctx) error {
	list, err := s.repo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": list})
}

// ====================================
// @Summary Get file by ID
// @Description Mendapatkan detail file berdasarkan ID
// @Tags File
// @Security BearerAuth
// @Param id path string true "ID File"
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Router /api/files/{id} [get]
func (s *FileService) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	f, err := s.repo.FindByID(oid)
	if err != nil || f == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "File tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": f})
}

// ====================================
// @Summary Hapus file
// @Description Hapus file berdasarkan ID (admin / pemilik file)
// @Tags File
// @Security BearerAuth
// @Param id path string true "ID File"
// @Success 200 {object} map[string]interface{}
// @Failure 400,401,403,404,500 {object} map[string]interface{}
// @Router /api/files/{id} [delete]
func (s *FileService) DeleteByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak valid"})
	}

	f, err := s.repo.FindByID(oid)
	if err != nil || f == nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "File tidak ditemukan"})
	}

	userAny := c.Locals("user")
	if userAny == nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Token tidak valid"})
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

	// Hanya admin atau uploader yang boleh hapus
	if role != "admin" && f.UploadedBy != userID {
		return c.Status(403).JSON(fiber.Map{"success": false, "message": "tidak punya akses menghapus file ini"})
	}

	os.Remove(f.FilePath)
	if err := s.repo.DeleteByID(oid); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "file dihapus"})
}

// ====================================
// helper simpan ke disk
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
