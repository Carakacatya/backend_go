package mocks

import (
	"context"
	"praktikum3/app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepositoryMock struct {
	FindByUsernameOrEmailFunc func(ctx context.Context, username string) (*model.User, error)
	SoftDeleteUserFunc        func(ctx context.Context, id primitive.ObjectID) error
}

func (m *UserRepositoryMock) FindByUsernameOrEmail(ctx context.Context, username string) (*model.User, error) {
	return m.FindByUsernameOrEmailFunc(ctx, username)
}

func (m *UserRepositoryMock) SoftDeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return m.SoftDeleteUserFunc(ctx, id)
}
