package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alumni struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NIM        string             `bson:"nim" json:"nim" example:"2020101234"`
	Nama       string             `bson:"nama" json:"nama" example:"Budi Santoso"`
	Jurusan    string             `bson:"jurusan" json:"jurusan" example:"Teknik Informatika"`
	Angkatan   int                `bson:"angkatan" json:"angkatan" example:"2020"`
	TahunLulus int                `bson:"tahun_lulus" json:"tahun_lulus" example:"2024"`
	Email      string             `bson:"email" json:"email" example:"budi@example.com"`
	NoTelepon  string             `bson:"no_telepon,omitempty" json:"no_telepon" example:"081234567890"`
	Alamat     string             `bson:"alamat,omitempty" json:"alamat" example:"Jl. Mawar No. 5"`
	CreatedAt  time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
	DeletedAt  *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id" example:"6710c5c2f8f4a385cd123456"`
}
