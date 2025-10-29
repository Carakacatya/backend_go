package repository

import (
	"context"
	"errors"
	"time"

	"praktikum3/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlumniRepository interface {
	GetAll() ([]model.Alumni, error)
	GetByID(id primitive.ObjectID) (*model.Alumni, error)
	Create(alumni *model.Alumni) error
	Update(id primitive.ObjectID, alumni *model.Alumni) error
	SoftDelete(id primitive.ObjectID) error
	GetTrashed() ([]model.Alumni, error)
	GetTrashedByID(id primitive.ObjectID) (*model.Alumni, error)
	Restore(id primitive.ObjectID) error
	ForceDelete(id primitive.ObjectID) error
}

type alumniRepository struct {
	col *mongo.Collection
}

func NewAlumniRepository(db *mongo.Database) AlumniRepository {
	return &alumniRepository{
		col: db.Collection("alumni"),
	}
}

func (r *alumniRepository) GetAll() ([]model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var alumniList []model.Alumni
	if err = cur.All(ctx, &alumniList); err != nil {
		return nil, err
	}
	return alumniList, nil
}

func (r *alumniRepository) GetByID(id primitive.ObjectID) (*model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var alumni model.Alumni
	err := r.col.FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&alumni)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &alumni, err
}

func (r *alumniRepository) Create(alumni *model.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	alumni.CreatedAt = time.Now()
	alumni.UpdatedAt = time.Now()
	_, err := r.col.InsertOne(ctx, alumni)
	return err
}

func (r *alumniRepository) Update(id primitive.ObjectID, alumni *model.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"nama":        alumni.Nama,
			"jurusan":     alumni.Jurusan,
			"angkatan":    alumni.Angkatan,
			"tahun_lulus": alumni.TahunLulus,
			"email":       alumni.Email,
			"no_telepon":  alumni.NoTelepon,
			"alamat":      alumni.Alamat,
			"updated_at":  time.Now(),
		},
	}

	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *alumniRepository) SoftDelete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *alumniRepository) GetTrashed() ([]model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var alumniList []model.Alumni
	if err = cur.All(ctx, &alumniList); err != nil {
		return nil, err
	}
	return alumniList, nil
}

func (r *alumniRepository) GetTrashedByID(id primitive.ObjectID) (*model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var alumni model.Alumni
	err := r.col.FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$ne": nil}}).Decode(&alumni)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &alumni, err
}

func (r *alumniRepository) Restore(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"deleted_at": nil,
			"updated_at": time.Now(),
		},
	}
	res, err := r.col.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$ne": nil}}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("data tidak ditemukan di trash")
	}
	return nil
}

func (r *alumniRepository) ForceDelete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("data tidak ditemukan")
	}
	return nil
}
