package repository

import (
	"context"
	"time"

	"praktikum3/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ✅ Interface (kontrak) untuk dipakai di layer service
type IUserRepository interface {
	FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error)
	SoftDeleteUser(ctx context.Context, id primitive.ObjectID) error
}

// ✅ Implementasi interface di bawah
type userRepository struct {
	col *mongo.Collection
}

// ✅ Constructor — mengembalikan interface (bukan struct langsung)
func NewUserRepository(db *mongo.Database) IUserRepository {
	return &userRepository{
		col: db.Collection("users"),
	}
}

// ✅ Cari user berdasarkan username atau email
func (r *userRepository) FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*model.User, error) {
	var user model.User
	err := r.col.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": usernameOrEmail},
			{"email": usernameOrEmail},
		},
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // tidak ada user ditemukan
		}
		return nil, err
	}

	return &user, nil
}

// ✅ Soft delete user berdasarkan ObjectID
func (r *userRepository) SoftDeleteUser(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		},
	}
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
