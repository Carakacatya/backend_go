package mocks

import (
	"praktikum3/app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PekerjaanRepositoryMock struct {
	GetAllWithQueryFunc   func(search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error)
	CountFunc             func(search string) (int, error)
	GetByIDFunc           func(id primitive.ObjectID) (*model.Pekerjaan, error)
	GetByAlumniIDFunc     func(alumniID primitive.ObjectID, includeDeleted bool) ([]model.Pekerjaan, error)
	CreateFunc            func(in model.CreatePekerjaanReq, mulai, selesai *time.Time) (primitive.ObjectID, error)
	UpdateFunc            func(id primitive.ObjectID, in model.UpdatePekerjaanReq, mulai, selesai *time.Time) error
	SoftDeleteByUserFunc  func(id primitive.ObjectID, alumniID primitive.ObjectID) error
	SoftDeleteByAdminFunc func(id primitive.ObjectID) error
	RestoreByIDFunc       func(id primitive.ObjectID) error
	HardDeleteByIDFunc    func(id primitive.ObjectID) error
	GetAllTrashFunc       func() ([]model.PekerjaanTrash, error)
	GetUserTrashFunc      func(alumniID primitive.ObjectID) ([]model.PekerjaanTrash, error)
}

// HardDeleteByUser implements repository.PekerjaanRepository.
func (m *PekerjaanRepositoryMock) HardDeleteByUser(id primitive.ObjectID, alumniID primitive.ObjectID) error {
	panic("unimplemented")
}

// RestoreByIDAndUser implements repository.PekerjaanRepository.
func (m *PekerjaanRepositoryMock) RestoreByIDAndUser(id primitive.ObjectID, alumniID primitive.ObjectID) error {
	panic("unimplemented")
}

func (m *PekerjaanRepositoryMock) GetAllWithQuery(search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
	if m.GetAllWithQueryFunc != nil {
		return m.GetAllWithQueryFunc(search, sortBy, order, limit, offset)
	}
	return nil, nil
}

func (m *PekerjaanRepositoryMock) Count(search string) (int, error) {
	if m.CountFunc != nil {
		return m.CountFunc(search)
	}
	return 0, nil
}

func (m *PekerjaanRepositoryMock) GetByID(id primitive.ObjectID) (*model.Pekerjaan, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *PekerjaanRepositoryMock) GetByAlumniID(alumniID primitive.ObjectID, includeDeleted bool) ([]model.Pekerjaan, error) {
	if m.GetByAlumniIDFunc != nil {
		return m.GetByAlumniIDFunc(alumniID, includeDeleted)
	}
	return nil, nil
}

func (m *PekerjaanRepositoryMock) Create(in model.CreatePekerjaanReq, mulai, selesai *time.Time) (primitive.ObjectID, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(in, mulai, selesai)
	}
	return primitive.NilObjectID, nil
}

func (m *PekerjaanRepositoryMock) Update(id primitive.ObjectID, in model.UpdatePekerjaanReq, mulai, selesai *time.Time) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, in, mulai, selesai)
	}
	return nil
}

func (m *PekerjaanRepositoryMock) SoftDeleteByUser(id primitive.ObjectID, alumniID primitive.ObjectID) error {
	if m.SoftDeleteByUserFunc != nil {
		return m.SoftDeleteByUserFunc(id, alumniID)
	}
	return nil
}

func (m *PekerjaanRepositoryMock) SoftDeleteByAdmin(id primitive.ObjectID) error {
	if m.SoftDeleteByAdminFunc != nil {
		return m.SoftDeleteByAdminFunc(id)
	}
	return nil
}

func (m *PekerjaanRepositoryMock) RestoreByID(id primitive.ObjectID) error {
	if m.RestoreByIDFunc != nil {
		return m.RestoreByIDFunc(id)
	}
	return nil
}

func (m *PekerjaanRepositoryMock) HardDeleteByID(id primitive.ObjectID) error {
	if m.HardDeleteByIDFunc != nil {
		return m.HardDeleteByIDFunc(id)
	}
	return nil
}

func (m *PekerjaanRepositoryMock) GetAllTrash() ([]model.PekerjaanTrash, error) {
	if m.GetAllTrashFunc != nil {
		return m.GetAllTrashFunc()
	}
	return nil, nil
}

func (m *PekerjaanRepositoryMock) GetUserTrash(alumniID primitive.ObjectID) ([]model.PekerjaanTrash, error) {
	if m.GetUserTrashFunc != nil {
		return m.GetUserTrashFunc(alumniID)
	}
	return nil, nil
}
