package repository

import (
	"praktikum3/app/model"
	"praktikum3/database"
)

func GetAlumniByStatus(status string) ([]model.AlumniPekerjaanReport, int, error) {
    rows, err := database.DB.Query(`
        SELECT a.id, a.nama, a.jurusan, a.angkatan,
               p.bidang_industri, p.nama_perusahaan, p.posisi_jabatan,
               p.tanggal_mulai_kerja, p.gaji_range,
               CASE 
                   WHEN p.tanggal_mulai_kerja <= NOW() - INTERVAL '1 year' 
                   THEN true ELSE false
               END AS lebih_dari_satu_tahun
        FROM alumni a
        JOIN pekerjaan_alumni p ON a.id = p.alumni_id
        WHERE 
            ($1 = 'aktif' AND p.status_pekerjaan = 'aktif')
            OR ($1 = 'tidak-aktif' AND p.status_pekerjaan IN ('selesai','resigned'))
    `, status)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var reports []model.AlumniPekerjaanReport
    for rows.Next() {
        var r model.AlumniPekerjaanReport
        if err := rows.Scan(
            &r.ID, &r.Nama, &r.Jurusan, &r.Angkatan,
            &r.BidangIndustri, &r.NamaPerusahaan, &r.PosisiJabatan,
            &r.TanggalMulaiKerja, &r.GajiRange, &r.LebihDariSatuTahun,
        ); err != nil {
            return nil, 0, err
        }
        reports = append(reports, r)
    }

    var count int
    err = database.DB.QueryRow(`
        SELECT COUNT(*)
        FROM pekerjaan_alumni
        WHERE 
            (
              ($1 = 'aktif' AND status_pekerjaan = 'aktif')
              OR ($1 = 'tidak-aktif' AND status_pekerjaan IN ('selesai','resigned'))
            )
          AND tanggal_mulai_kerja <= NOW() - INTERVAL '1 year'
    `, status).Scan(&count)
    if err != nil {
        return nil, 0, err
    }

    return reports, count, nil
}
