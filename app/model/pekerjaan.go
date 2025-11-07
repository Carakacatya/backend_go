package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Pekerjaan struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AlumniID            primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	NamaPerusahaan      string             `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string             `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string             `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string             `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           *string            `bson:"gaji_range,omitempty" json:"gaji_range,omitempty"`
	TanggalMulaiKerja   *time.Time         `bson:"tanggal_mulai_kerja,omitempty" json:"tanggal_mulai_kerja,omitempty"`
	TanggalSelesaiKerja *time.Time         `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string             `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  *string            `bson:"deskripsi_pekerjaan,omitempty" json:"deskripsi_pekerjaan,omitempty"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt           *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
type PekerjaanTrash struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AlumniID            primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	NamaPerusahaan      string             `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string             `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string             `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string             `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           *string            `bson:"gaji_range,omitempty" json:"gaji_range,omitempty"`
	TanggalMulaiKerja   *time.Time         `bson:"tanggal_mulai_kerja,omitempty" json:"tanggal_mulai_kerja,omitempty"`
	TanggalSelesaiKerja *time.Time         `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string             `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  *string            `bson:"deskripsi_pekerjaan,omitempty" json:"deskripsi_pekerjaan,omitempty"`
	DeletedAt           *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
