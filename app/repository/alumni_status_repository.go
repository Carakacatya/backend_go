package repository

import (
	"context"
	"log"
	"time"

	"praktikum3/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlumniStatusRepository interface {
	GetAlumniByStatus(status string) ([]model.AlumniPekerjaanReport, int, error)
}

type alumniStatusRepository struct {
	alumniCol    *mongo.Collection
	pekerjaanCol *mongo.Collection
}

func NewAlumniStatusRepository(db *mongo.Database) AlumniStatusRepository {
	return &alumniStatusRepository{
		alumniCol:    db.Collection("alumni"),
		pekerjaanCol: db.Collection("pekerjaan_alumni"),
	}
}

func (r *alumniStatusRepository) GetAlumniByStatus(status string) ([]model.AlumniPekerjaanReport, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Filter pekerjaan berdasarkan status
	var pekerjaanFilter bson.M
	if status == "aktif" {
		pekerjaanFilter = bson.M{"status_pekerjaan": "aktif"}
	} else if status == "tidak-aktif" {
		pekerjaanFilter = bson.M{
			"status_pekerjaan": bson.M{"$in": []string{"selesai", "resigned"}},
		}
	} else {
		pekerjaanFilter = bson.M{} // semua
	}

	cursor, err := r.pekerjaanCol.Find(ctx, pekerjaanFilter)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reports []model.AlumniPekerjaanReport
	count := 0
	now := time.Now()

	for cursor.Next(ctx) {
		var p struct {
			AlumniID          interface{} `bson:"alumni_id"`
			NamaPerusahaan    string      `bson:"nama_perusahaan"`
			PosisiJabatan     string      `bson:"posisi_jabatan"`
			BidangIndustri    string      `bson:"bidang_industri"`
			TanggalMulaiKerja time.Time   `bson:"tanggal_mulai_kerja"`
			GajiRange         string      `bson:"gaji_range"`
		}
		if err := cursor.Decode(&p); err != nil {
			log.Println("decode error:", err)
			continue
		}

		// cari data alumni terkait
		var a struct {
			ID        interface{} `bson:"_id"`
			Nama      string      `bson:"nama"`
			Jurusan   string      `bson:"jurusan"`
			Angkatan  int         `bson:"angkatan"`
		}

		err := r.alumniCol.FindOne(ctx, bson.M{"_id": p.AlumniID}).Decode(&a)
		if err != nil {
			log.Printf("alumni %v tidak ditemukan: %v", p.AlumniID, err)
			continue
		}

		// logika lama: lebih dari 1 tahun bekerja
		lebihSetahun := p.TanggalMulaiKerja.Before(now.AddDate(-1, 0, 0))
		if lebihSetahun {
			count++
		}

		reports = append(reports, model.AlumniPekerjaanReport{
			Nama:               a.Nama,
			Jurusan:            a.Jurusan,
			Angkatan:           a.Angkatan,
			BidangIndustri:     p.BidangIndustri,
			NamaPerusahaan:     p.NamaPerusahaan,
			PosisiJabatan:      p.PosisiJabatan,
			TanggalMulaiKerja:  p.TanggalMulaiKerja,
			GajiRange:          p.GajiRange,
			LebihDariSatuTahun: lebihSetahun,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return reports, count, nil
}
