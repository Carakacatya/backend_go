package service

import (
	"errors"
	"praktikum3/app/repository"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) SoftDeleteUser(id uint, role string) error {
	if role != "admin" {
		return errors.New("forbidden: only admin can delete user")
	}

	return s.UserRepo.SoftDeleteUser(id)
}
