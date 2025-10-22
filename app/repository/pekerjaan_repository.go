package repository

import (
	"database/sql"
	"fmt"
	"praktikum3/app/model"
	"strings"
	"time"
)

type PekerjaanRepository struct {
	DB *sql.DB
}

func NewPekerjaanRepository(db *sql.DB) *PekerjaanRepository {
	return &PekerjaanRepository{DB: db}
}

func (r *PekerjaanRepository) GetAllTrash() ([]model.PekerjaanTrash, error) {
	rows, err := r.DB.Query(`
		SELECT id, nama_perusahaan, posisi_jabatan, deleted_at, updated_at
		FROM pekerjaan_alumni
		WHERE deleted_at IS NOT NULL;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.PekerjaanTrash
	for rows.Next() {
		var p model.PekerjaanTrash
		if err := rows.Scan(
			&p.ID,
			&p.NamaPerusahaan,
			&p.PosisiJabatan,
			&p.DeletedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	// Kalau tidak ada hasil, kembalikan slice kosong, bukan error
	if len(result) == 0 {
		return []model.PekerjaanTrash{}, nil
	}

	return result, nil
}

func (r *PekerjaanRepository) GetAllWithQuery(search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	allowedSort := map[string]string{
		"id": "p.id", "alumni_id": "p.alumni_id", "nama_perusahaan": "p.nama_perusahaan",
		"posisi_jabatan": "p.posisi_jabatan", "bidang_industri": "p.bidang_industri",
		"lokasi_kerja": "p.lokasi_kerja", "tanggal_mulai_kerja": "p.tanggal_mulai_kerja",
		"status_pekerjaan": "p.status_pekerjaan", "created_at": "p.created_at",
	}
	col, ok := allowedSort[strings.ToLower(sortBy)]
	if !ok {
		col = "p.created_at"
	}
	if strings.ToUpper(order) != "ASC" {
		order = "DESC"
	}

	query := fmt.Sprintf(`
		SELECT p.id, p.alumni_id, a.nama, p.nama_perusahaan, p.posisi_jabatan,
		       p.bidang_industri, p.lokasi_kerja, p.gaji_range,
		       p.tanggal_mulai_kerja, p.tanggal_selesai_kerja,
		       p.status_pekerjaan, p.deskripsi_pekerjaan,
		       p.created_at, p.updated_at, p.deleted_at
		FROM pekerjaan_alumni p
		LEFT JOIN alumni a ON p.alumni_id = a.id
		WHERE (a.nama ILIKE $1 OR p.nama_perusahaan ILIKE $1 
		       OR p.posisi_jabatan ILIKE $1 OR p.bidang_industri ILIKE $1)
		  AND p.deleted_at IS NULL
		ORDER BY %s %s
		LIMIT $2 OFFSET $3`, col, order)

	rows, err := r.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		err := rows.Scan(
			&p.ID, &p.AlumniID, &p.AlumniNama, &p.NamaPerusahaan,
			&p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
			&p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PekerjaanRepository) Count(search string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM pekerjaan_alumni p
		LEFT JOIN alumni a ON p.alumni_id = a.id
		WHERE (a.nama ILIKE $1 OR p.nama_perusahaan ILIKE $1 
		       OR p.posisi_jabatan ILIKE $1 OR p.bidang_industri ILIKE $1)
		  AND p.deleted_at IS NULL
	`, "%"+search+"%").Scan(&total)
	return total, err
}

func (r *PekerjaanRepository) GetByID(id int) (*model.Pekerjaan, error) {
	var p model.Pekerjaan
	err := r.DB.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, deleted_at
		FROM pekerjaan_alumni 
		WHERE id = $1
	`, id).Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
		&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
		&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PekerjaanRepository) GetByAlumniID(alumniID int) ([]model.Pekerjaan, error) {
	rows, err := r.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, deleted_at
		FROM pekerjaan_alumni 
		WHERE alumni_id = $1 AND deleted_at IS NULL
	`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PekerjaanRepository) Create(in model.CreatePekerjaanReq, mulai, selesai *time.Time) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		INSERT INTO pekerjaan_alumni (
			alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
			lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
			status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())
		RETURNING id
	`,
		in.AlumniID, in.NamaPerusahaan, in.PosisiJabatan, in.BidangIndustri,
		in.LokasiKerja, in.GajiRange, mulai, selesai,
		in.StatusPekerjaan, in.DeskripsiPekerjaan,
	).Scan(&id)
	return id, err
}

func (r *PekerjaanRepository) Update(id int, in model.UpdatePekerjaanReq, mulai, selesai *time.Time) error {
	_, err := r.DB.Exec(`
		UPDATE pekerjaan_alumni 
		SET nama_perusahaan=$1, posisi_jabatan=$2, bidang_industri=$3,
		    lokasi_kerja=$4, gaji_range=$5, tanggal_mulai_kerja=$6, tanggal_selesai_kerja=$7,
		    status_pekerjaan=$8, deskripsi_pekerjaan=$9, updated_at=NOW()
		WHERE id=$10 AND deleted_at IS NULL
	`,
		in.NamaPerusahaan, in.PosisiJabatan, in.BidangIndustri,
		in.LokasiKerja, in.GajiRange, mulai, selesai,
		in.StatusPekerjaan, in.DeskripsiPekerjaan, id,
	)
	return err
}

// SoftDeleteByUser hanya bisa digunakan oleh user biasa.
// Ia hanya boleh menghapus data pekerjaan miliknya sendiri.
func (r *PekerjaanRepository) SoftDeleteByUser(id int, userID int) error {
	query := `
		UPDATE pekerjaan_alumni 
		SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND alumni_id IN (SELECT id FROM alumni WHERE user_id = $2)
		AND deleted_at IS NULL
		RETURNING deleted_at
	`

	var deletedAt time.Time
	err := r.DB.QueryRow(query, id, userID).Scan(&deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("gagal soft delete: data tidak ditemukan atau bukan milik user")
		}
		return err
	}

	return nil
}

// SoftDeleteByAdmin digunakan oleh admin untuk menghapus data siapa pun.
func (r *PekerjaanRepository) SoftDeleteByAdmin(id int) error {
	query := `
		UPDATE pekerjaan_alumni 
		SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
		RETURNING deleted_at
	`

	var deletedAt time.Time
	err := r.DB.QueryRow(query, id).Scan(&deletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("gagal soft delete: data tidak ditemukan")
		}
		return err
	}

	return nil
}

func (r *PekerjaanRepository) RestoreByID(id int) error {
	_, err := r.DB.Exec(`
		UPDATE pekerjaan_alumni 
		SET deleted_at = NULL, updated_at = NOW() 
		WHERE id = $1
	`, id)
	return err
}

func (r *PekerjaanRepository) RestoreByIDAndUser(id, userID int) error {
	_, err := r.DB.Exec(`
		UPDATE pekerjaan_alumni
		SET deleted_at = NULL, updated_at = NOW()
		WHERE id = $1 AND alumni_id IN (SELECT id FROM alumni WHERE user_id = $2)
	`, id, userID)
	return err
}

func (r *PekerjaanRepository) HardDeleteByID(id int) error {
	_, err := r.DB.Exec(`DELETE FROM pekerjaan_alumni WHERE id = $1`, id)
	return err
}

func (r *PekerjaanRepository) HardDeleteByUser(id, userID int) error {
	_, err := r.DB.Exec(`
		DELETE FROM pekerjaan_alumni
		WHERE id = $1 AND alumni_id IN (SELECT id FROM alumni WHERE user_id = $2)
	`, id, userID)
	return err
}

func (r *PekerjaanRepository) GetUserTrash(userID int) ([]model.PekerjaanTrash, error) {
	rows, err := r.DB.Query(`
		SELECT p.id, p.nama_perusahaan, p.posisi_jabatan, p.deskripsi, p.alumni_id, p.deleted_at, p.updated_at
		FROM pekerjaan_alumni p
		JOIN alumni a ON p.alumni_id = a.id
		WHERE a.user_id = $1 AND p.deleted_at IS NOT NULL
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.PekerjaanTrash
	for rows.Next() {
		var p model.PekerjaanTrash
		if err := rows.Scan(
			&p.ID,
			&p.NamaPerusahaan,
			&p.PosisiJabatan,
			&p.DeletedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	// Jika tidak ada data, bisa return slice kosong tanpa error
	return result, nil
}

