package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File merepresentasikan metadata file yang diupload
type File struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	AlumniID     *primitive.ObjectID `bson:"alumni_id,omitempty" json:"alumni_id,omitempty"` // bisa nil jika general
	FileName     string              `bson:"file_name" json:"file_name"`
	OriginalName string              `bson:"original_name" json:"original_name"`
	FilePath     string              `bson:"file_path" json:"file_path"`
	FileSize     int64               `bson:"file_size" json:"file_size"`
	FileType     string              `bson:"file_type" json:"file_type"`
	Category     string              `bson:"category" json:"category"` // "photo" atau "certificate"
	UploadedBy   primitive.ObjectID  `bson:"uploaded_by" json:"uploaded_by"`
	UploadedAt   time.Time           `bson:"uploaded_at" json:"uploaded_at"`
}
