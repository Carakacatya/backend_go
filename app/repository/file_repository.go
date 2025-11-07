package repository

import (
	"context"
	"time"

	"praktikum3/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository interface {
	Create(file *model.File) (primitive.ObjectID, error)
	FindAll() ([]model.File, error)
	FindByID(id primitive.ObjectID) (*model.File, error)
	DeleteByID(id primitive.ObjectID) error
	FindByUploadedBy(userID primitive.ObjectID) ([]model.File, error)
}

type fileRepository struct {
	col *mongo.Collection
}

func NewFileRepository(db *mongo.Database) FileRepository {
	return &fileRepository{
		col: db.Collection("files"),
	}
}

// ✅ Insert file metadata ke MongoDB
func (r *fileRepository) Create(file *model.File) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	file.UploadedAt = time.Now()

	res, err := r.col.InsertOne(ctx, file)
	if err != nil {
		return primitive.NilObjectID, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, mongo.ErrNilDocument
	}

	return oid, nil
}

// ✅ Ambil semua file
func (r *fileRepository) FindAll() ([]model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.File
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// ✅ Ambil file berdasarkan ID
func (r *fileRepository) FindByID(id primitive.ObjectID) (*model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var f model.File
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// ✅ Hapus file berdasarkan ID
func (r *fileRepository) DeleteByID(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ✅ Ambil semua file berdasarkan uploader
func (r *fileRepository) FindByUploadedBy(userID primitive.ObjectID) ([]model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.col.Find(ctx, bson.M{"uploaded_by": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.File
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}
