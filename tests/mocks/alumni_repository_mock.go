package mocks

import (
	"praktikum3/app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlumniRepositoryMock struct {
	// GetAll
	GetAllFunc func() ([]model.Alumni, error)

	// GetByID
	GetByIDFunc func(id primitive.ObjectID) (*model.Alumni, error)

	// Create
	CreateFunc func(alumni *model.Alumni) error

	// Update
	UpdateFunc func(id primitive.ObjectID, alumni *model.Alumni) error

	// SoftDelete
	SoftDeleteFunc func(id primitive.ObjectID) error

	// GetTrashed
	GetTrashedFunc func() ([]model.Alumni, error)

	// GetTrashedByID
	GetTrashedByIDFunc func(id primitive.ObjectID) (*model.Alumni, error)

	// Restore
	RestoreFunc func(id primitive.ObjectID) error

	// ForceDelete
	ForceDeleteFunc func(id primitive.ObjectID) error
}

func (m *AlumniRepositoryMock) GetAll() ([]model.Alumni, error) {
    if m.GetAllFunc != nil {
        return m.GetAllFunc()
    }
    return []model.Alumni{}, nil
}

func (m *AlumniRepositoryMock) GetByID(id primitive.ObjectID) (*model.Alumni, error) {
    if m.GetByIDFunc != nil {
        return m.GetByIDFunc(id)
    }
    return nil, nil
}

func (m *AlumniRepositoryMock) Create(alumni *model.Alumni) error {
    if m.CreateFunc != nil {
        return m.CreateFunc(alumni)
    }
    return nil
}

func (m *AlumniRepositoryMock) Update(id primitive.ObjectID, alumni *model.Alumni) error {
    if m.UpdateFunc != nil {
        return m.UpdateFunc(id, alumni)
    }
    return nil
}

func (m *AlumniRepositoryMock) SoftDelete(id primitive.ObjectID) error {
    if m.SoftDeleteFunc != nil {
        return m.SoftDeleteFunc(id)
    }
    return nil
}

func (m *AlumniRepositoryMock) GetTrashed() ([]model.Alumni, error) {
    if m.GetTrashedFunc != nil {
        return m.GetTrashedFunc()
    }
    return []model.Alumni{}, nil
}

func (m *AlumniRepositoryMock) GetTrashedByID(id primitive.ObjectID) (*model.Alumni, error) {
    if m.GetTrashedByIDFunc != nil {
        return m.GetTrashedByIDFunc(id)
    }
    return nil, nil
}

func (m *AlumniRepositoryMock) Restore(id primitive.ObjectID) error {
    if m.RestoreFunc != nil {
        return m.RestoreFunc(id)
    }
    return nil
}

func (m *AlumniRepositoryMock) ForceDelete(id primitive.ObjectID) error {
    if m.ForceDeleteFunc != nil {
        return m.ForceDeleteFunc(id)
    }
    return nil
}

