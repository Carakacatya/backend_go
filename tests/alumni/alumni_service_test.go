package alumni_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"praktikum3/app/model"
	"praktikum3/app/service"
	"praktikum3/tests/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestApp(repo *mocks.AlumniRepositoryMock) *fiber.App {
    app := fiber.New()

    alumniService := service.NewAlumniService(repo)

    app.Use(func(c *fiber.Ctx) error {
        c.Locals("user", map[string]interface{}{
            "id": primitive.NewObjectID().Hex(),
            "username": "test",
            "role": "admin",
        })
        c.Locals("role", "admin")
        return c.Next()
    })

    // ORDER FIXED !!!
    app.Get("/alumni", alumniService.GetAll)

    app.Get("/alumni/trash", alumniService.GetTrashed)
    app.Put("/alumni/restore/:id", alumniService.Restore)
    app.Delete("/alumni/hard/:id", alumniService.HardDelete)

    app.Get("/alumni/:id", alumniService.GetByID)
    app.Post("/alumni", alumniService.Create)
    app.Put("/alumni/:id", alumniService.Update)
    app.Delete("/alumni/:id", alumniService.SoftDelete)

    return app
}


//
// ================================
// TEST GET ALL
// ================================
func TestGetAll_Error(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetAllFunc: func() ([]model.Alumni, error) {
			return nil, errors.New("db error")
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/alumni", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetAll_Success(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetAllFunc: func() ([]model.Alumni, error) {
			return []model.Alumni{}, nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/alumni", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

//
// ================================
// TEST GET BY ID
// ================================
func TestGetByID_InvalidID(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("GET", "/alumni/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetByIDFunc: func(id primitive.ObjectID) (*model.Alumni, error) {
			return nil, nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/alumni/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestGetByID_Success(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetByIDFunc: func(id primitive.ObjectID) (*model.Alumni, error) {
			return &model.Alumni{Nama: "Test"}, nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/alumni/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

//
// ================================
// TEST CREATE
// ================================
func TestCreate_InvalidBody(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("POST", "/alumni", bytes.NewBufferString("INVALID"))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCreate_MissingFields(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{}
	app := setupTestApp(repo)

	body, _ := json.Marshal(model.Alumni{
		Email: "mail@mail.com",
	})
	req := httptest.NewRequest("POST", "/alumni", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCreate_RepoError(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		CreateFunc: func(a *model.Alumni) error {
			return errors.New("insert error")
		},
	}

	app := setupTestApp(repo)

	body, _ := json.Marshal(model.Alumni{
		Nama:  "Test",
		Email: "mail@mail.com",
	})

	req := httptest.NewRequest("POST", "/alumni", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestCreate_Success(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		CreateFunc: func(a *model.Alumni) error {
			return nil
		},
	}

	app := setupTestApp(repo)

	body, _ := json.Marshal(model.Alumni{
		Nama:  "Test",
		Email: "mail@mail.com",
	})

	req := httptest.NewRequest("POST", "/alumni", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

//
// ================================
// TEST UPDATE
// ================================
func TestUpdate_InvalidID(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("PUT", "/alumni/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestUpdate_InvalidBody(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("PUT", "/alumni/"+primitive.NewObjectID().Hex(), bytes.NewBufferString("INVALID"))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestUpdate_NotFound(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetByIDFunc: func(id primitive.ObjectID) (*model.Alumni, error) {
			return nil, nil
		},
	}

	app := setupTestApp(repo)

	body, _ := json.Marshal(model.Alumni{Nama: "New"})
	req := httptest.NewRequest("PUT", "/alumni/"+primitive.NewObjectID().Hex(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestUpdate_RepoError(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetByIDFunc: func(id primitive.ObjectID) (*model.Alumni, error) {
			return &model.Alumni{}, nil
		},
		UpdateFunc: func(id primitive.ObjectID, a *model.Alumni) error {
			return errors.New("update error")
		},
	}

	app := setupTestApp(repo)

	body, _ := json.Marshal(model.Alumni{Nama: "New"})
	req := httptest.NewRequest("PUT", "/alumni/"+primitive.NewObjectID().Hex(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestUpdate_Success(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetByIDFunc: func(id primitive.ObjectID) (*model.Alumni, error) {
			return &model.Alumni{}, nil
		},
		UpdateFunc: func(id primitive.ObjectID, a *model.Alumni) error {
			return nil
		},
	}

	app := setupTestApp(repo)

	body, _ := json.Marshal(model.Alumni{Nama: "New"})
	req := httptest.NewRequest("PUT", "/alumni/"+primitive.NewObjectID().Hex(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

//
// ================================
// TEST SOFT DELETE
// ================================
func TestSoftDelete_InvalidID(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("DELETE", "/alumni/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestSoftDelete_RepoError(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		SoftDeleteFunc: func(id primitive.ObjectID) error {
			return errors.New("delete error")
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("DELETE", "/alumni/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestSoftDelete_Success(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		SoftDeleteFunc: func(id primitive.ObjectID) error {
			return nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("DELETE", "/alumni/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

//
// ================================
// TEST GET TRASHED
// ================================
func TestGetTrashed_Error(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetTrashedFunc: func() ([]model.Alumni, error) {
			return nil, errors.New("trash error")
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/alumni/trash", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetTrashed_Success(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetTrashedFunc: func() ([]model.Alumni, error) {
			return []model.Alumni{}, nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("GET", "/alumni/trash", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

//
// ================================
// TEST RESTORE
// ================================
func TestRestore_InvalidID(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("PUT", "/alumni/restore/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestRestore_NotFound(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetTrashedByIDFunc: func(id primitive.ObjectID) (*model.Alumni, error) {
			return nil, nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("PUT", "/alumni/restore/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestRestore_RepoError(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetTrashedByIDFunc: func(id primitive.ObjectID) (*model.Alumni, error) {
			return &model.Alumni{}, nil
		},
		RestoreFunc: func(id primitive.ObjectID) error {
			return errors.New("restore error")
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("PUT", "/alumni/restore/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestRestore_Success(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetTrashedByIDFunc: func(id primitive.ObjectID) (*model.Alumni, error) {
			return &model.Alumni{}, nil
		},
		RestoreFunc: func(id primitive.ObjectID) error {
			return nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("PUT", "/alumni/restore/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

//
// ================================
// TEST HARD DELETE
// ================================
func TestHardDelete_InvalidID(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{}
	app := setupTestApp(repo)

	req := httptest.NewRequest("DELETE", "/alumni/hard/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestHardDelete_NotFound(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetTrashedByIDFunc: func(id primitive.ObjectID) (*model.Alumni, error) {
			return nil, nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("DELETE", "/alumni/hard/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestHardDelete_RepoError(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetTrashedByIDFunc: func(id primitive.ObjectID) (*model.Alumni, error) {
			return &model.Alumni{}, nil
		},
		ForceDeleteFunc: func(id primitive.ObjectID) error {
			return errors.New("delete error")
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("DELETE", "/alumni/hard/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestHardDelete_Success(t *testing.T) {
	repo := &mocks.AlumniRepositoryMock{
		GetTrashedByIDFunc: func(id primitive.ObjectID) (*model.Alumni, error) {
			return &model.Alumni{}, nil
		},
		ForceDeleteFunc: func(id primitive.ObjectID) error {
			return nil
		},
	}

	app := setupTestApp(repo)
	req := httptest.NewRequest("DELETE", "/alumni/hard/"+primitive.NewObjectID().Hex(), nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}
