package mocks

import "praktikum3/app/model"

type TokenGeneratorMock struct {
	GenerateFunc func(user model.User) (string, error)
}

func (m TokenGeneratorMock) Generate(user model.User) (string, error) {
	return m.GenerateFunc(user)
}
