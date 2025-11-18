package mocks

import (
	"praktikum3/app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileRepositoryMock struct {
	CreateFunc           func(file *model.File) (primitive.ObjectID, error)
	FindAllFunc          func() ([]model.File, error)
	FindByIDFunc         func(id primitive.ObjectID) (*model.File, error)
	DeleteByIDFunc       func(id primitive.ObjectID) error
	FindByUploadedByFunc func(userID primitive.ObjectID) ([]model.File, error)
}

func (m *FileRepositoryMock) Create(file *model.File) (primitive.ObjectID, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(file)
	}
	return primitive.NilObjectID, nil
}

func (m *FileRepositoryMock) FindAll() ([]model.File, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc()
	}
	return nil, nil
}

func (m *FileRepositoryMock) FindByID(id primitive.ObjectID) (*model.File, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, nil
}

func (m *FileRepositoryMock) DeleteByID(id primitive.ObjectID) error {
	if m.DeleteByIDFunc != nil {
		return m.DeleteByIDFunc(id)
	}
	return nil
}

func (m *FileRepositoryMock) FindByUploadedBy(userID primitive.ObjectID) ([]model.File, error) {
	if m.FindByUploadedByFunc != nil {
		return m.FindByUploadedByFunc(userID)
	}
	return nil, nil
}
