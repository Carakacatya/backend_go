package service

import (
	"context"
	"errors"
	"praktikum3/app/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ✅ Struct service untuk mengelola operasi user
type UserService struct {
	userRepo repository.IUserRepository
}

// ✅ Constructor: menerima dependency dari repository layer
func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{userRepo: repo}
}

// ✅ SoftDeleteUser — hanya boleh dijalankan oleh admin
func (s *UserService) SoftDeleteUser(ctx context.Context, id string, role string) error {
	// Cek role user (hanya admin)
	if role != "admin" {
		return errors.New("forbidden: hanya admin yang dapat menghapus user")
	}

	// Konversi string ke ObjectID MongoDB
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID user tidak valid")
	}

	// Jalankan operasi soft delete di repository
	return s.userRepo.SoftDeleteUser(ctx, objectID)
}
