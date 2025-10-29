package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// ================== CREATE REQUEST ==================
type CreatePekerjaanReq struct {
	AlumniID            primitive.ObjectID `json:"alumni_id" bson:"alumni_id"`               // ✅ langsung ObjectID
	NamaPerusahaan      string             `json:"nama_perusahaan" bson:"nama_perusahaan"`
	PosisiJabatan       string             `json:"posisi_jabatan" bson:"posisi_jabatan"`
	BidangIndustri      string             `json:"bidang_industri" bson:"bidang_industri"`
	LokasiKerja         string             `json:"lokasi_kerja" bson:"lokasi_kerja"`
	GajiRange           *string            `json:"gaji_range,omitempty" bson:"gaji_range,omitempty"`
	TanggalMulaiKerja   string             `json:"tanggal_mulai_kerja,omitempty" bson:"tanggal_mulai_kerja,omitempty"`
	TanggalSelesaiKerja string             `json:"tanggal_selesai_kerja,omitempty" bson:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string             `json:"status_pekerjaan" bson:"status_pekerjaan"`
	DeskripsiPekerjaan  *string            `json:"deskripsi_pekerjaan,omitempty" bson:"deskripsi_pekerjaan,omitempty"`
}

// ================== UPDATE REQUEST ==================
type UpdatePekerjaanReq struct {
	NamaPerusahaan      string  `json:"nama_perusahaan" bson:"nama_perusahaan"`
	PosisiJabatan       string  `json:"posisi_jabatan" bson:"posisi_jabatan"`
	BidangIndustri      string  `json:"bidang_industri" bson:"bidang_industri"`
	LokasiKerja         string  `json:"lokasi_kerja" bson:"lokasi_kerja"`
	GajiRange           *string `json:"gaji_range,omitempty" bson:"gaji_range,omitempty"`
	TanggalMulaiKerja   string  `json:"tanggal_mulai_kerja,omitempty" bson:"tanggal_mulai_kerja,omitempty"`
	TanggalSelesaiKerja string  `json:"tanggal_selesai_kerja,omitempty" bson:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string  `json:"status_pekerjaan" bson:"status_pekerjaan"`
	DeskripsiPekerjaan  *string `json:"deskripsi_pekerjaan,omitempty" bson:"deskripsi_pekerjaan,omitempty"`
}
