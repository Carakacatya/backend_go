package repository

import (
	"database/sql"
	"errors"
	"praktikum3/app/model"
)

type AlumniRepository interface {
	GetAll() ([]model.Alumni, error)
	GetByID(id int) (*model.Alumni, error)
	Create(alumni *model.Alumni) error
	Update(alumni *model.Alumni) error
	SoftDelete(id int) error
	GetTrashed() ([]model.Alumni, error)
	GetTrashedByID(id int) (*model.Alumni, error)
	Restore(id int) error
	ForceDelete(id int) error
}

type alumniRepository struct {
	db *sql.DB
}

func NewAlumniRepository(db *sql.DB) AlumniRepository {
	return &alumniRepository{db: db}
}

func (r *alumniRepository) GetAll() ([]model.Alumni, error) {
	rows, err := r.db.Query(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, user_id, deleted_at FROM alumni WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		var a model.Alumni
		err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.UserID, &a.DeletedAt)
		if err != nil {
			return nil, err
		}
		alumniList = append(alumniList, a)
	}
	return alumniList, nil
}

func (r *alumniRepository) GetByID(id int) (*model.Alumni, error) {
	row := r.db.QueryRow(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, user_id, deleted_at FROM alumni WHERE id = ? AND deleted_at IS NULL`, id)
	var a model.Alumni
	err := row.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.UserID, &a.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *alumniRepository) Create(alumni *model.Alumni) error {
	_, err := r.db.Exec(`INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		alumni.NIM, alumni.Nama, alumni.Jurusan, alumni.Angkatan, alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat, alumni.UserID)
	return err
}

func (r *alumniRepository) Update(alumni *model.Alumni) error {
	_, err := r.db.Exec(`UPDATE alumni SET nama = ?, jurusan = ?, angkatan = ?, tahun_lulus = ?, email = ?, no_telepon = ?, alamat = ? WHERE id = ?`,
		alumni.Nama, alumni.Jurusan, alumni.Angkatan, alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat, alumni.ID)
	return err
}

func (r *alumniRepository) SoftDelete(id int) error {
	_, err := r.db.Exec(`UPDATE alumni SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	return err
}

func (r *alumniRepository) GetTrashed() ([]model.Alumni, error) {
	rows, err := r.db.Query(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, user_id, deleted_at FROM alumni WHERE deleted_at IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		var a model.Alumni
		err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.UserID, &a.DeletedAt)
		if err != nil {
			return nil, err
		}
		alumniList = append(alumniList, a)
	}
	return alumniList, nil
}

func (r *alumniRepository) GetTrashedByID(id int) (*model.Alumni, error) {
	row := r.db.QueryRow(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, user_id, deleted_at FROM alumni WHERE id = ? AND deleted_at IS NOT NULL`, id)
	var a model.Alumni
	err := row.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.UserID, &a.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *alumniRepository) Restore(id int) error {
	res, err := r.db.Exec(`UPDATE alumni SET deleted_at = NULL WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data tidak ditemukan di trash")
	}
	return nil
}

func (r *alumniRepository) ForceDelete(id int) error {
	res, err := r.db.Exec(`DELETE FROM alumni WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("data tidak ditemukan")
	}
	return nil
}
