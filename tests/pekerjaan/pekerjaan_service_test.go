package pekerjaan_test

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"praktikum3/app/model"
	"praktikum3/app/service"
	"praktikum3/tests/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ========================== SETUP APP ==========================
func setupTestApp(repo *mocks.PekerjaanRepositoryMock) *fiber.App {
	app := fiber.New()

	s := service.NewPekerjaanService(repo)

	// Bypass auth middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{
			"id":   primitive.NewObjectID().Hex(),
			"role": "admin", // default admin for testing
		})
		return c.Next()
	})

	// Route without real middleware
	app.Get("/pekerjaan", s.GetAll)
	app.Get("/pekerjaan/:id", s.GetByID)
	app.Get("/pekerjaan/alumni/:alumni_id", s.GetByAlumniID)
	app.Post("/pekerjaan", s.Create)
	app.Put("/pekerjaan/:id", s.Update)
	app.Delete("/pekerjaan/:id", s.SoftDelete)
	app.Put("/pekerjaan/restore/:id", s.Restore)
	app.Delete("/pekerjaan/hard/:id", s.HardDelete)
	app.Get("/pekerjaan/trash", s.GetTrashed)

	return app
}

// ==============================================================
//                       GET ALL
// ==============================================================
func TestGetAll_RepoError(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		GetAllWithQueryFunc: func(s1, s2, s3 string, l, o int) ([]model.Pekerjaan, error) {
			return nil, errors.New("db error")
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/pekerjaan", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetAll_CountError(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		GetAllWithQueryFunc: func(s1, s2, s3 string, l, o int) ([]model.Pekerjaan, error) {
			return []model.Pekerjaan{}, nil
		},
		CountFunc: func(s string) (int, error) {
			return 0, errors.New("count error")
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/pekerjaan", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetAll_Success(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		GetAllWithQueryFunc: func(s1, s2, s3 string, l, o int) ([]model.Pekerjaan, error) {
			return []model.Pekerjaan{}, nil
		},
		CountFunc: func(s string) (int, error) {
			return 10, nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/pekerjaan", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// ==============================================================
//                        GET BY ID
// ==============================================================
func TestGetByID_InvalidID(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("GET", "/pekerjaan/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestGetByID_RepoError(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		GetByIDFunc: func(id primitive.ObjectID) (*model.Pekerjaan, error) {
			return nil, errors.New("db error")
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/pekerjaan/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		GetByIDFunc: func(id primitive.ObjectID) (*model.Pekerjaan, error) {
			return nil, nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/pekerjaan/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestGetByID_Success(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		GetByIDFunc: func(id primitive.ObjectID) (*model.Pekerjaan, error) {
			return &model.Pekerjaan{NamaPerusahaan: "TestCorp"}, nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/pekerjaan/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// ==============================================================
//                     GET BY ALUMNI ID
// ==============================================================
func TestGetByAlumniID_InvalidID(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("GET", "/pekerjaan/alumni/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestGetByAlumniID_RepoError(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		GetByAlumniIDFunc: func(a primitive.ObjectID, b bool) ([]model.Pekerjaan, error) {
			return nil, errors.New("db error")
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/pekerjaan/alumni/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetByAlumniID_Success(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		GetByAlumniIDFunc: func(a primitive.ObjectID, b bool) ([]model.Pekerjaan, error) {
			return []model.Pekerjaan{}, nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/pekerjaan/alumni/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// ==============================================================
//                         CREATE
// ==============================================================
func TestCreate_InvalidBody(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("POST", "/pekerjaan", bytes.NewBufferString("INVALID"))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCreate_InvalidStartDate(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{}
	app := setupTestApp(repo)

	body := `{"alumni_id":"` + primitive.NewObjectID().Hex() + `","nama_perusahaan":"X","tanggal_mulai_kerja":"INVALID"}`
	req := httptest.NewRequest("POST", "/pekerjaan", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCreate_InvalidEndDate(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{}
	app := setupTestApp(repo)

	body := `{"alumni_id":"` + primitive.NewObjectID().Hex() + `","nama_perusahaan":"X","tanggal_mulai_kerja":"2020-01-01","tanggal_selesai_kerja":"WRONG"}`
	req := httptest.NewRequest("POST", "/pekerjaan", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCreate_RepoError(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		CreateFunc: func(in model.CreatePekerjaanReq, a, b *time.Time) (primitive.ObjectID, error) {
			return primitive.NilObjectID, errors.New("insert error")
		},
	}

	app := setupTestApp(repo)

	body := `{
		"alumni_id":"` + primitive.NewObjectID().Hex() + `",
		"nama_perusahaan":"X",
		"tanggal_mulai_kerja":"2020-01-01"
	}`

	req := httptest.NewRequest("POST", "/pekerjaan", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestCreate_Success(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		CreateFunc: func(in model.CreatePekerjaanReq, a, b *time.Time) (primitive.ObjectID, error) {
			return primitive.NewObjectID(), nil
		},
	}

	app := setupTestApp(repo)

	body := `{
		"alumni_id":"` + primitive.NewObjectID().Hex() + `",
		"nama_perusahaan":"X",
		"tanggal_mulai_kerja":"2020-01-01"
	}`

	req := httptest.NewRequest("POST", "/pekerjaan", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// ==============================================================
//                         UPDATE
// ==============================================================
func TestUpdate_InvalidID(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("PUT", "/pekerjaan/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestUpdate_InvalidBody(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("PUT", "/pekerjaan/"+primitive.NewObjectID().Hex(), bytes.NewBufferString("INVALID"))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUpdate_InvalidStartDate(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{}
	app := setupTestApp(repo)

	body := `{"tanggal_mulai_kerja":"INVALID"}`
	req := httptest.NewRequest("PUT", "/pekerjaan/"+primitive.NewObjectID().Hex(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestUpdate_RepoError(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		UpdateFunc: func(id primitive.ObjectID, in model.UpdatePekerjaanReq, a, b *time.Time) error {
			return errors.New("update error")
		},
	}

	app := setupTestApp(repo)

	body := `{
		"tanggal_mulai_kerja":"2020-01-01"
	}`

	req := httptest.NewRequest("PUT", "/pekerjaan/"+primitive.NewObjectID().Hex(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestUpdate_Success(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		UpdateFunc: func(id primitive.ObjectID, in model.UpdatePekerjaanReq, a, b *time.Time) error {
			return nil
		},
	}

	app := setupTestApp(repo)

	body := `{
		"tanggal_mulai_kerja":"2020-01-01"
	}`

	req := httptest.NewRequest("PUT", "/pekerjaan/"+primitive.NewObjectID().Hex(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

// ==============================================================
//                     SOFT DELETE
// ==============================================================
func TestSoftDelete_InvalidID(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("DELETE", "/pekerjaan/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestSoftDelete_AdminError(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		SoftDeleteByAdminFunc: func(id primitive.ObjectID) error {
			return errors.New("delete error")
		},
	}

	app := setupTestApp(repo)

	req := httptest.NewRequest("DELETE", "/pekerjaan/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestSoftDelete_AdminSuccess(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		SoftDeleteByAdminFunc: func(id primitive.ObjectID) error {
			return nil
		},
	}

	app := setupTestApp(repo)

	req := httptest.NewRequest("DELETE", "/pekerjaan/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// ==============================================================
//                         RESTORE
// ==============================================================
func TestRestore_InvalidID(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("PUT", "/pekerjaan/restore/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestRestore_RepoError(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		RestoreByIDFunc: func(id primitive.ObjectID) error {
			return errors.New("restore error")
		},
	}

	app := setupTestApp(repo)

	req := httptest.NewRequest("PUT", "/pekerjaan/restore/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestRestore_Success(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		RestoreByIDFunc: func(id primitive.ObjectID) error {
			return nil
		},
	}

	app := setupTestApp(repo)

	req := httptest.NewRequest("PUT", "/pekerjaan/restore/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// ==============================================================
//                      HARD DELETE
// ==============================================================
func TestHardDelete_InvalidID(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("DELETE", "/pekerjaan/hard/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestHardDelete_RepoError(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		HardDeleteByIDFunc: func(id primitive.ObjectID) error {
			return errors.New("delete error")
		},
	}

	app := setupTestApp(repo)

	req := httptest.NewRequest("DELETE", "/pekerjaan/hard/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestHardDelete_Success(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		HardDeleteByIDFunc: func(id primitive.ObjectID) error {
			return nil
		},
	}

	app := setupTestApp(repo)

	req := httptest.NewRequest("DELETE", "/pekerjaan/hard/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

// ==============================================================
//                      GET TRASH
// ==============================================================

func TestGetTrashed_Error(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		GetAllTrashFunc: func() ([]model.PekerjaanTrash, error) {
			return nil, errors.New("trash error")
		},
	}

	app := fiber.New()
	s := service.NewPekerjaanService(repo)

	// Inject user admin (wajib, karena trash adalah admin-only)
	app.Get("/pekerjaan/trash", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{
			"id":   primitive.NewObjectID().Hex(),
			"role": "admin",
		})
		return s.GetTrashed(c)
	})

	req := httptest.NewRequest("GET", "/pekerjaan/trash", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetTrashed_Success(t *testing.T) {
	repo := &mocks.PekerjaanRepositoryMock{
		GetAllTrashFunc: func() ([]model.PekerjaanTrash, error) {
			return []model.PekerjaanTrash{}, nil
		},
	}

	app := fiber.New()
	s := service.NewPekerjaanService(repo)

	// Inject user admin (wajib supaya tidak 400)
	app.Get("/pekerjaan/trash", func(c *fiber.Ctx) error {
		c.Locals("user", map[string]interface{}{
			"id":   primitive.NewObjectID().Hex(),
			"role": "admin",
		})
		return s.GetTrashed(c)
	})

	req := httptest.NewRequest("GET", "/pekerjaan/trash", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

