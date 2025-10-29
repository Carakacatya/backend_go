package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"praktikum3/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PekerjaanRepository interface {
	GetAllWithQuery(search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error)
	Count(search string) (int, error)
	GetByID(id primitive.ObjectID) (*model.Pekerjaan, error)
	GetByAlumniID(alumniID primitive.ObjectID, includeDeleted bool) ([]model.Pekerjaan, error)
	Create(in model.CreatePekerjaanReq, mulai, selesai *time.Time) (primitive.ObjectID, error)
	Update(id primitive.ObjectID, in model.UpdatePekerjaanReq, mulai, selesai *time.Time) error
	SoftDeleteByUser(id primitive.ObjectID, alumniID primitive.ObjectID) error
	SoftDeleteByAdmin(id primitive.ObjectID) error
	RestoreByID(id primitive.ObjectID) error
	RestoreByIDAndUser(id primitive.ObjectID, alumniID primitive.ObjectID) error
	HardDeleteByID(id primitive.ObjectID) error
	HardDeleteByUser(id primitive.ObjectID, alumniID primitive.ObjectID) error
	GetAllTrash() ([]model.PekerjaanTrash, error)
	GetUserTrash(alumniID primitive.ObjectID) ([]model.PekerjaanTrash, error)
}

type pekerjaanRepository struct {
	col *mongo.Collection
}

func NewPekerjaanRepository(db *mongo.Database) PekerjaanRepository {
	return &pekerjaanRepository{
		col: db.Collection("pekerjaan_alumni"),
	}
}

// ================= GET ALL =================
func (r *pekerjaanRepository) GetAllWithQuery(search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"deleted_at": nil,
		"$or": []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
		},
	}

	opts := options.Find()
	if order == "ASC" {
		opts.SetSort(bson.D{{Key: sortBy, Value: 1}})
	} else {
		opts.SetSort(bson.D{{Key: sortBy, Value: -1}})
	}
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.Pekerjaan
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// ================= COUNT =================
func (r *pekerjaanRepository) Count(search string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"deleted_at": nil,
		"$or": []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
		},
	}

	count, err := r.col.CountDocuments(ctx, filter)
	return int(count), err
}

// ================= GET BY ID =================
func (r *pekerjaanRepository) GetByID(id primitive.ObjectID) (*model.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var p model.Pekerjaan
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &p, err
}

// ================= GET BY ALUMNI ID =================
func (r *pekerjaanRepository) GetByAlumniID(alumniID primitive.ObjectID, includeDeleted bool) ([]model.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"alumni_id": alumniID},
			{"alumni_id": alumniID.Hex()},
		},
	}
	if !includeDeleted {
		filter["deleted_at"] = nil
	}

	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Pekerjaan
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// ================= CREATE =================
func (r *pekerjaanRepository) Create(in model.CreatePekerjaanReq, mulai, selesai *time.Time) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := bson.M{
		"alumni_id":             in.AlumniID,
		"nama_perusahaan":       in.NamaPerusahaan,
		"posisi_jabatan":        in.PosisiJabatan,
		"bidang_industri":       in.BidangIndustri,
		"lokasi_kerja":          in.LokasiKerja,
		"gaji_range":            in.GajiRange,
		"tanggal_mulai_kerja":   mulai,
		"tanggal_selesai_kerja": selesai,
		"status_pekerjaan":      in.StatusPekerjaan,
		"deskripsi_pekerjaan":   in.DeskripsiPekerjaan,
		"created_at":            time.Now(),
		"updated_at":            time.Now(),
	}

	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("gagal konversi ObjectID")
	}
	return oid, nil
}

// ================= UPDATE =================
func (r *pekerjaanRepository) Update(id primitive.ObjectID, in model.UpdatePekerjaanReq, mulai, selesai *time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{
		"nama_perusahaan":       in.NamaPerusahaan,
		"posisi_jabatan":        in.PosisiJabatan,
		"bidang_industri":       in.BidangIndustri,
		"lokasi_kerja":          in.LokasiKerja,
		"gaji_range":            in.GajiRange,
		"tanggal_mulai_kerja":   mulai,
		"tanggal_selesai_kerja": selesai,
		"status_pekerjaan":      in.StatusPekerjaan,
		"deskripsi_pekerjaan":   in.DeskripsiPekerjaan,
		"updated_at":            time.Now(),
	}}

	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// ================= SOFT DELETE =================
func (r *pekerjaanRepository) SoftDeleteByUser(id primitive.ObjectID, alumniID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
		"$or": []bson.M{
			{"alumni_id": alumniID},
			{"alumni_id": alumniID.Hex()},
		},
		"deleted_at": nil,
	}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now(), "updated_at": time.Now()}}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("data tidak ditemukan atau bukan milik user")
	}
	return nil
}

func (r *pekerjaanRepository) SoftDeleteByAdmin(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"deleted_at": time.Now(), "updated_at": time.Now()}}
	res, err := r.col.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": nil}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("data tidak ditemukan")
	}
	return nil
}

// ================= RESTORE =================
func (r *pekerjaanRepository) RestoreByID(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.M{"$set": bson.M{"deleted_at": nil, "updated_at": time.Now()}}
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *pekerjaanRepository) RestoreByIDAndUser(id, alumniID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{
		"_id": id,
		"$or": []bson.M{
			{"alumni_id": alumniID},
			{"alumni_id": alumniID.Hex()},
		},
	}
	update := bson.M{"$set": bson.M{"deleted_at": nil, "updated_at": time.Now()}}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

// ================= HARD DELETE =================
func (r *pekerjaanRepository) HardDeleteByID(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *pekerjaanRepository) HardDeleteByUser(id, alumniID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.col.DeleteOne(ctx, bson.M{
		"_id": id,
		"$or": []bson.M{
			{"alumni_id": alumniID},
			{"alumni_id": alumniID.Hex()},
		},
	})
	return err
}

// ================= TRASH =================
func (r *pekerjaanRepository) GetAllTrash() ([]model.PekerjaanTrash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := r.col.Find(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var trash []model.PekerjaanTrash
	if err := cur.All(ctx, &trash); err != nil {
		return nil, err
	}
	return trash, nil
}

func (r *pekerjaanRepository) GetUserTrash(alumniID primitive.ObjectID) ([]model.PekerjaanTrash, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"alumni_id": alumniID},
			{"alumni_id": alumniID.Hex()},
		},
		"deleted_at": bson.M{"$ne": nil},
	}
	cur, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var trash []model.PekerjaanTrash
	if err := cur.All(ctx, &trash); err != nil {
		return nil, err
	}
	return trash, nil
}
