package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlumniPekerjaanReport struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nama              string             `bson:"nama" json:"nama"`
	Jurusan           string             `bson:"jurusan" json:"jurusan"`
	Angkatan          int                `bson:"angkatan" json:"angkatan"`
	BidangIndustri    string             `bson:"bidang_industri" json:"bidang_industri"`
	NamaPerusahaan    string             `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan     string             `bson:"posisi_jabatan" json:"posisi_jabatan"`
	TanggalMulaiKerja time.Time          `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	GajiRange         string             `bson:"gaji_range" json:"gaji_range"`
	LebihDariSatuTahun bool              `bson:"lebih_dari_satu_tahun" json:"lebih_dari_satu_tahun"`
}
