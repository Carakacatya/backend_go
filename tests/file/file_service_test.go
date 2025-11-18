package file_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"praktikum3/app/model"
	"praktikum3/app/service"
	"praktikum3/tests/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// helper: create multipart request with a file part (allow custom content-type)
func makeMultipartRequest(t *testing.T, method, url, fieldName, filename, contentType string, data []byte, formFields map[string]string) *http.Request {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// create part with header so we can set Content-Type
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="`+fieldName+`"; filename="`+filename+`"`)
	if contentType != "" {
		h.Set("Content-Type", contentType)
	} else {
		h.Set("Content-Type", "application/octet-stream")
	}
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatalf("createpart err: %v", err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatalf("write part err: %v", err)
	}

	for k, v := range formFields {
		if err := w.WriteField(k, v); err != nil {
			t.Fatalf("write field err: %v", err)
		}
	}

	w.Close()

	req := httptest.NewRequest(method, url, &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

// ====================== setup app helper ======================
func setupTestApp(repo *mocks.FileRepositoryMock, uploadBase string) *fiber.App {
	app := fiber.New()
	svc := service.NewFileService(repo, uploadBase)

	// Bypass auth middleware to inject user (by default admin)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{
			"id":   primitive.NewObjectID().Hex(),
			"role": "admin",
		})
		return c.Next()
	})

	// routes (no middleware)
	app.Post("/api/files/photo", svc.UploadPhoto)
	app.Post("/api/files/certificate", svc.UploadCertificate)
	app.Get("/api/files", svc.GetAll)
	app.Get("/api/files/:id", svc.GetByID)
	app.Delete("/api/files/:id", svc.DeleteByID)

	return app
}

// ==================== UPLOAD PHOTO ====================
func TestUploadPhoto_NoFile(t *testing.T) {
	repo := &mocks.FileRepositoryMock{}
	app := setupTestApp(repo, t.TempDir())

	req := httptest.NewRequest("POST", "/api/files/photo", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUploadPhoto_InvalidType(t *testing.T) {
	repo := &mocks.FileRepositoryMock{}
	app := setupTestApp(repo, t.TempDir())

	data := []byte("plain text")
	req := makeMultipartRequest(t, "POST", "/api/files/photo", "file", "test.txt", "text/plain", data, nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUploadPhoto_MaxSizeExceeded(t *testing.T) {
	repo := &mocks.FileRepositoryMock{}
	app := setupTestApp(repo, t.TempDir())

	// create >1MB payload
	data := make([]byte, 1*1024*1024+10)
	req := makeMultipartRequest(t, "POST", "/api/files/photo", "file", "big.jpg", "image/jpeg", data, nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUploadPhoto_AdminWithoutAlumniID(t *testing.T) {
	repo := &mocks.FileRepositoryMock{}
	app := setupTestApp(repo, t.TempDir())

	data := []byte{0xFF, 0xD8, 0xFF} // jpeg magic bytes (small)
	req := makeMultipartRequest(t, "POST", "/api/files/photo", "file", "photo.jpg", "image/jpeg", data, nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUploadPhoto_AdminInvalidAlumniID(t *testing.T) {
	repo := &mocks.FileRepositoryMock{}
	app := setupTestApp(repo, t.TempDir())

	data := []byte{0xFF, 0xD8, 0xFF}
	req := makeMultipartRequest(t, "POST", "/api/files/photo", "file", "photo.jpg", "image/jpeg", data, map[string]string{"alumni_id": "INVALID"})
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUploadPhoto_UserSetAlumniID_Forbidden(t *testing.T) {
	// change middleware to non-admin
	repo := &mocks.FileRepositoryMock{}
	uploadBase := t.TempDir()
	app := fiber.New()
	svc := service.NewFileService(repo, uploadBase)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"id": primitive.NewObjectID().Hex(), "role": "user"})
		return c.Next()
	})
	app.Post("/api/files/photo", svc.UploadPhoto)

	data := []byte{0xFF, 0xD8, 0xFF}
	req := makeMultipartRequest(t, "POST", "/api/files/photo", "file", "photo.jpg", "image/jpeg", data, map[string]string{"alumni_id": primitive.NewObjectID().Hex()})
	resp, _ := app.Test(req)
	assert.Equal(t, 403, resp.StatusCode)
}

func TestUploadPhoto_SaveFileError(t *testing.T) {
	repo := &mocks.FileRepositoryMock{}
	uploadBase := t.TempDir()

	// create a file with the same name as "photos" folder so MkdirAll will fail
	badDir := filepath.Join(uploadBase, "photos")
	if err := os.WriteFile(badDir, []byte("not a dir"), 0644); err != nil {
		t.Fatalf("setup err: %v", err)
	}

	app := setupTestApp(repo, uploadBase)

	data := []byte{0xFF, 0xD8, 0xFF}
	req := makeMultipartRequest(t, "POST", "/api/files/photo", "file", "photo.jpg", "image/jpeg", data, map[string]string{"alumni_id": primitive.NewObjectID().Hex()})
	resp, _ := app.Test(req)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestUploadPhoto_RepoError_Rollback(t *testing.T) {
	repo := &mocks.FileRepositoryMock{
		CreateFunc: func(f *model.File) (primitive.ObjectID, error) {
			return primitive.NilObjectID, errors.New("insert error")
		},
	}
	uploadBase := t.TempDir()
	app := setupTestApp(repo, uploadBase)

	data := []byte{0xFF, 0xD8, 0xFF}
	req := makeMultipartRequest(t, "POST", "/api/files/photo", "file", "photo.jpg", "image/jpeg", data, map[string]string{"alumni_id": primitive.NewObjectID().Hex()})
	resp, _ := app.Test(req)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestUploadPhoto_Success(t *testing.T) {
	repo := &mocks.FileRepositoryMock{
		CreateFunc: func(f *model.File) (primitive.ObjectID, error) {
			return primitive.NewObjectID(), nil
		},
	}
	uploadBase := t.TempDir()
	app := setupTestApp(repo, uploadBase)

	data := []byte{0xFF, 0xD8, 0xFF}
	req := makeMultipartRequest(t, "POST", "/api/files/photo", "file", "photo.jpg", "image/jpeg", data, map[string]string{"alumni_id": primitive.NewObjectID().Hex()})
	resp, _ := app.Test(req)
	assert.Equal(t, 201, resp.StatusCode)
}

// ==================== UPLOAD CERTIFICATE (similar tests) ====================
func TestUploadCertificate_InvalidType(t *testing.T) {
	repo := &mocks.FileRepositoryMock{}
	app := setupTestApp(repo, t.TempDir())

	data := []byte("%PDF-") // not actually pdf but content-type matters
	req := makeMultipartRequest(t, "POST", "/api/files/certificate", "file", "file.txt", "text/plain", data, map[string]string{"alumni_id": primitive.NewObjectID().Hex()})
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUploadCertificate_MaxSizeExceeded(t *testing.T) {
	repo := &mocks.FileRepositoryMock{}
	app := setupTestApp(repo, t.TempDir())

	data := make([]byte, 2*1024*1024+10)
	req := makeMultipartRequest(t, "POST", "/api/files/certificate", "file", "cert.pdf", "application/pdf", data, map[string]string{"alumni_id": primitive.NewObjectID().Hex()})
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUploadCertificate_Success(t *testing.T) {
	repo := &mocks.FileRepositoryMock{
		CreateFunc: func(f *model.File) (primitive.ObjectID, error) {
			return primitive.NewObjectID(), nil
		},
	}
	app := setupTestApp(repo, t.TempDir())

	data := []byte("%PDF-1.4")
	req := makeMultipartRequest(t, "POST", "/api/files/certificate", "file", "cert.pdf", "application/pdf", data, map[string]string{"alumni_id": primitive.NewObjectID().Hex()})
	resp, _ := app.Test(req)
	assert.Equal(t, 201, resp.StatusCode)
}

// ==================== GET ALL ====================
func TestGetAll_Error(t *testing.T) {
	repo := &mocks.FileRepositoryMock{
		FindAllFunc: func() ([]model.File, error) {
			return nil, errors.New("db error")
		},
	}
	app := setupTestApp(repo, t.TempDir())
	req := httptest.NewRequest("GET", "/api/files", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetAll_Success(t *testing.T) {
	repo := &mocks.FileRepositoryMock{
		FindAllFunc: func() ([]model.File, error) {
			return []model.File{}, nil
		},
	}
	app := setupTestApp(repo, t.TempDir())
	req := httptest.NewRequest("GET", "/api/files", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

// ==================== GET BY ID ====================
func TestGetByID_InvalidID(t *testing.T) {
	repo := &mocks.FileRepositoryMock{}
	app := setupTestApp(repo, t.TempDir())
	req := httptest.NewRequest("GET", "/api/files/INVALID", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mocks.FileRepositoryMock{
		FindByIDFunc: func(id primitive.ObjectID) (*model.File, error) {
			return nil, errors.New("not found")
		},
	}
	app := setupTestApp(repo, t.TempDir())
	req := httptest.NewRequest("GET", "/api/files/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestGetByID_Success(t *testing.T) {
	repo := &mocks.FileRepositoryMock{
		FindByIDFunc: func(id primitive.ObjectID) (*model.File, error) {
			return &model.File{FileName: "a.jpg"}, nil
		},
	}
	app := setupTestApp(repo, t.TempDir())
	req := httptest.NewRequest("GET", "/api/files/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

// ==================== DELETE ====================
func TestDelete_InvalidID(t *testing.T) {
	repo := &mocks.FileRepositoryMock{}
	app := setupTestApp(repo, t.TempDir())
	req := httptest.NewRequest("DELETE", "/api/files/INVALID", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestDelete_NotFound(t *testing.T) {
	repo := &mocks.FileRepositoryMock{
		FindByIDFunc: func(id primitive.ObjectID) (*model.File, error) {
			return nil, errors.New("not found")
		},
	}
	app := setupTestApp(repo, t.TempDir())
	req := httptest.NewRequest("DELETE", "/api/files/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestDelete_Unauthorized(t *testing.T) {
	repo := &mocks.FileRepositoryMock{
		FindByIDFunc: func(id primitive.ObjectID) (*model.File, error) {
			return &model.File{FileName: "a.jpg", UploadedBy: primitive.NewObjectID()}, nil
		},
	}
	// setup app without user
	app := fiber.New()
	svc := service.NewFileService(repo, t.TempDir())
	app.Delete("/api/files/:id", svc.DeleteByID)
	req := httptest.NewRequest("DELETE", "/api/files/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestDelete_Forbidden(t *testing.T) {
	// file exists uploaded by other user, current role user
	uploader := primitive.NewObjectID()
	repo := &mocks.FileRepositoryMock{
		FindByIDFunc: func(id primitive.ObjectID) (*model.File, error) {
			return &model.File{FileName: "a.jpg", UploadedBy: uploader}, nil
		},
	}
	app := fiber.New()
	svc := service.NewFileService(repo, t.TempDir())
	// set user as non-admin and different id
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"id": primitive.NewObjectID().Hex(), "role": "user"})
		return c.Next()
	})
	app.Delete("/api/files/:id", svc.DeleteByID)
	req := httptest.NewRequest("DELETE", "/api/files/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 403, resp.StatusCode)
}

func TestDelete_RepoError(t *testing.T) {
	oid := primitive.NewObjectID()
	repo := &mocks.FileRepositoryMock{
		FindByIDFunc: func(id primitive.ObjectID) (*model.File, error) {
			return &model.File{FileName: "a.jpg", UploadedBy: oid}, nil
		},
		DeleteByIDFunc: func(id primitive.ObjectID) error {
			return errors.New("delete db error")
		},
	}
	app := fiber.New()
	svc := service.NewFileService(repo, t.TempDir())
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"id": oid.Hex(), "role": "admin"})
		return c.Next()
	})
	app.Delete("/api/files/:id", svc.DeleteByID)
	req := httptest.NewRequest("DELETE", "/api/files/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestDelete_Success(t *testing.T) {
	oid := primitive.NewObjectID()
	// To avoid error when removing file path, ensure path exists
	td := t.TempDir()
	filePath := filepath.Join(td, "f.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0644); err != nil {
		t.Fatalf("write tmp file: %v", err)
	}

	repo := &mocks.FileRepositoryMock{
		FindByIDFunc: func(id primitive.ObjectID) (*model.File, error) {
			return &model.File{FileName: "f.txt", UploadedBy: oid, FilePath: filePath}, nil
		},
		DeleteByIDFunc: func(id primitive.ObjectID) error {
			return nil
		},
	}
	app := fiber.New()
	svc := service.NewFileService(repo, td)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{"id": oid.Hex(), "role": "admin"})
		return c.Next()
	})
	app.Delete("/api/files/:id", svc.DeleteByID)
	req := httptest.NewRequest("DELETE", "/api/files/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}
