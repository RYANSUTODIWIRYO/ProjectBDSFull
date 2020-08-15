package services

import (
	ent "bds/entities"
	"database/sql"
)

type Service struct {
	DB *sql.DB
}

func (s Service) LoginUser(id_user int, password string) (ent.User, error) {
	rows, err := s.DB.Query("SELECT * FROM user WHERE id_user = ? AND password = ? ", id_user, password)
	if err != nil {
		panic(err)
	} else {
		var user ent.User
		for rows.Next() {
			var id_user int64
			var password string
			var nama string
			var role string
			var cabang string
			err2 := rows.Scan(&id_user, &password, &nama, &role, &cabang)
			if err2 != nil {
				panic(err)
			}
			user = ent.User{
				Id_user:   id_user,
				Password:  password,
				Nama_user: nama,
				Role:      role,
				Cabang:    cabang,
			}
		}
		return user, nil
	}
}

func (s Service) CariNasabah(rekeningTujuan int) (ent.NasabahDetail, error) {
	//rows, err := us.DB.Query("SELECT * FROM user WHERE id_user = ? AND password = ? ",id_user, password)
	rows, err := s.DB.Query(
		"SELECT nasabah_detail.cif, nasabah.nama, nasabah_detail.no_rekening, nasabah_detail.saldo "+
			"FROM bank.nasabah_detail "+
			"INNER JOIN bank.nasabah "+
			"ON (nasabah_detail.cif = nasabah.cif AND nasabah_detail.no_rekening = ?)", rekeningTujuan)

	if err != nil {
		panic(err)
	} else {
		var nasabahDetail ent.NasabahDetail
		for rows.Next() {
			var (
				cif         int64
				nama        string
				no_rekening int64
				saldo       int64
			)
			err2 := rows.Scan(&cif, &nama, &no_rekening, &saldo)
			if err2 != nil {
				panic(err)
			}
			nasabahDetail = ent.NasabahDetail{
				Cif:         cif,
				Nama:        nama,
				No_rekening: no_rekening,
				Saldo:       saldo,
			}
		}
		return nasabahDetail, nil
	}
}

// insert setor tunai to db
func (s Service) SetorTunaiService(transaksi ent.Transaksi) (int, ent.Transaksi, error) {
	transaksi.Saldo = transaksi.Saldo + transaksi.Nominal
	rows, err := s.DB.Exec(
		"INSERT INTO transaksi (id_user, no_rekening, tanggal, jenis_transaksi, nominal, saldo, berita) values (?,?,?,?,?,?,?)",
		transaksi.Id_user,
		transaksi.No_rekening,
		transaksi.Tanggal,
		transaksi.Jenis_transaksi,
		transaksi.Nominal,
		transaksi.Saldo,
		transaksi.Berita)
	setor, err := s.DB.Exec(
		"UPDATE nasabah_detail SET saldo = ? where no_rekening = ?", transaksi.Saldo, transaksi.No_rekening)
	if err != nil {
		panic(err)
	}
	update1, _ := setor.RowsAffected()
	update2, _ := rows.RowsAffected()
	transaksi.Saldo = transaksi.Saldo
	if update1 > 0 && update2 > 0 {
		return int(update1), transaksi, nil
	} else {
		return 0, ent.Transaksi{}, nil
	}
}

// insert setor tunai to db
func (s Service) TarikTunaiService(transaksi ent.Transaksi, nasabah ent.NasabahDetail) (int, ent.Transaksi, error) {

	nasabah.Saldo = nasabah.Saldo - transaksi.Nominal
	if nasabah.Saldo < 0 {
		return -1, ent.Transaksi{}, nil
	}

	rows, err := s.DB.Exec(
		"INSERT INTO transaksi (id_user, no_rekening, tanggal, jenis_transaksi, nominal, saldo, berita) values (?,?,?,?,?,?,?)",
		transaksi.Id_user,
		nasabah.No_rekening,
		transaksi.Tanggal,
		transaksi.Jenis_transaksi,
		transaksi.Nominal,
		nasabah.Saldo,
		transaksi.Berita)
	_, err = s.DB.Exec(
		"UPDATE nasabah_detail SET saldo = ? where no_rekening = ?", nasabah.Saldo, nasabah.No_rekening)
	if err != nil {
		panic(err)
	} else {
		transaksi.Saldo = nasabah.Saldo
		rows.RowsAffected()
		return 1, transaksi, nil
	}
}

func (s Service) CetakBuku(no_rekening int) ([]ent.Transaksi, error) {
	rows, err := s.DB.Query("SELECT * FROM transaksi WHERE no_rekening = ?", no_rekening)

	if err != nil {
		return []ent.Transaksi{}, err
	} else {
		var transaksi []ent.Transaksi
		for rows.Next() {
			var id_transaksi int64
			var id_user int64
			var no_rekening int64
			var tanggal string
			var jenis_transaksi string
			var nominal int64
			var saldo int64
			var berita string
			err2 := rows.Scan(&id_transaksi, &id_user, &no_rekening, &tanggal, &jenis_transaksi, &nominal, &saldo, &berita)
			if err2 != nil {
				return []ent.Transaksi{}, err2
			} else {
				trx := ent.Transaksi{
					Id_transaksi:    id_transaksi,
					Id_user:         id_user,
					No_rekening:     no_rekening,
					Tanggal:         tanggal,
					Jenis_transaksi: jenis_transaksi,
					Nominal:         nominal,
					Saldo:           saldo,
					Berita:          berita,
				}
				transaksi = append(transaksi, trx)
			}
		}
		return transaksi, nil
	}
}

func (s Service) PindahBukuService(idUser int64, tanggal string, rekeningAwal, rekeingTujuan ent.NasabahDetail, nominal int64, berita string) (int64, error) {
	// check saldo overload
	if rekeningAwal.Saldo < nominal {
		return 0, nil
	} else {
		saldoRekAwal := int(rekeningAwal.Saldo - nominal)
		//fmt.Println(saldoRekAwal)
		saldoRekTujuan := int(rekeingTujuan.Saldo + nominal)
		//fmt.Println(saldoRekTujuan)
		_, err := s.DB.Exec(
			"INSERT INTO transaksi (id_user, no_rekening, tanggal, jenis_transaksi, nominal, saldo, berita) VALUES (?,?,?,?,?,?,?)",
			idUser, rekeningAwal.No_rekening, tanggal, "pb (d)", nominal, saldoRekAwal, berita)
		if err != nil {
			panic(err)
		}
		//fmt.Println(insertRekUtama)
		_, err2 := s.DB.Exec("UPDATE nasabah_detail SET saldo = ? WHERE no_rekening = ?", saldoRekAwal, rekeningAwal.No_rekening)
		if err2 != nil {
			panic(err)
		}
		//fmt.Println(update)
		_, err3 := s.DB.Exec("INSERT INTO transaksi (id_user, no_rekening, tanggal, jenis_transaksi, nominal, saldo, berita) VALUES (?,?,?,?,?,?,?)",
			idUser, rekeingTujuan.No_rekening, tanggal, "pb (k)", nominal, saldoRekTujuan, berita)
		if err3 != nil {
			panic(err)
		}
		//fmt.Println(insertRekKedua)
		_, err4 := s.DB.Exec("UPDATE nasabah_detail SET saldo = ? where no_rekening = ?", saldoRekTujuan, rekeingTujuan.No_rekening)
		if err4 != nil {
			panic(err)
		}
		//fmt.Println(updateRekTujuan)

	}
	return idUser, nil
}

func (s Service) FindByCifOrNikService(cif int64) (ent.Nasabah, error) {
	//Mencari nasabah berdasarkan cif atau nik
	rows, err := s.DB.Query("SELECT * FROM nasabah WHERE cif = ? OR nik = ?", cif, cif)
	if err != nil {
		panic(err)
	} else {
		var nasabah ent.Nasabah
		for rows.Next() {
			var (
				cif           int64
				nik           int64
				nama          string
				tempat_lahir  string
				tanggal_lahir string
				alamat        string
				no_telepon    string
			)

			//Menampung hasil query
			err2 := rows.Scan(&cif, &nik, &nama, &tempat_lahir, &tanggal_lahir, &alamat, &no_telepon)
			if err2 != nil {
				panic(err2)
			}

			nasabah = ent.Nasabah{
				Cif:           cif,
				Nik:           nik,
				Nama:          nama,
				Tempat_lahir:  tempat_lahir,
				Tanggal_lahir: tanggal_lahir,
				Alamat:        alamat,
				No_telepon:    no_telepon,
			}
			//fmt.Println("called")
			//fmt.Println(nasabah)
		}
		return nasabah, nil
	}
}

func (s Service) BuatCifService(nasabah ent.Nasabah) (ent.Nasabah, error) {
	//Memasukan data ke table nasabah
	rows, err := s.DB.Exec("INSERT INTO nasabah (nik,nama,tempat_lahir,tanggal_lahir,alamat,no_telepon) values (?,?,?,?,?,?)",
		nasabah.Nik, nasabah.Nama, nasabah.Tempat_lahir, nasabah.Tanggal_lahir, nasabah.Alamat, nasabah.No_telepon)
	//fmt.Println("nik", nasabah.Nik)
	if err != nil {
		panic(err)
	} else {
		status, _ := rows.RowsAffected()
		if status > 0 {
			response, err := s.FindByCifOrNikService(nasabah.Nik)
			nasabah.Cif = response.Cif
			return nasabah, err
		} else {
			return ent.Nasabah{}, err
		}
	}
}

func (s Service) FindLastRekService() (int64, error) {
	rows, err := s.DB.Query("SELECT no_rekening FROM nasabah_detail ORDER BY no_rekening DESC LIMIT 1")
	if err != nil {
		panic(err)
	} else {
		// looping data
		var no_rekening int64
		for rows.Next() {

			err2 := rows.Scan(&no_rekening)
			if err2 != nil {
				panic(err2)
			}
		}
		return no_rekening, nil
	}
}

func (s Service) BuatTabunganService(nasabah ent.NasabahDetail) (ent.NasabahDetail, error) {
	//Mencari nomor rekening yang dimasukan terakhir
	last_no_rekening, _ := s.FindLastRekService()
	last_no_rekening += 1

	//Memasukan ke database nasabah_detail
	//fmt.Println(last_no_rekening)
	rows, err := s.DB.Exec("INSERT INTO nasabah_detail (cif, no_rekening, saldo) VALUES (?,?,?)",
		nasabah.Cif, last_no_rekening, nasabah.Saldo)
	if err != nil {
		panic(err)
	}

	//Mencari nama dari method FindByCifOrNikService()
	status, _ := rows.RowsAffected()
	if status > 0 {
		// call method print nasabah by cif
		response, _ := s.FindByCifOrNikService(nasabah.Cif)
		nasabah.Nama = response.Nama
		nasabah.No_rekening = last_no_rekening
		return nasabah, nil
	} else {
		return ent.NasabahDetail{}, nil
	}
}

func (s Service) UpdateNasabahService(nasabah ent.Nasabah) (ent.Nasabah, error) {

	rows, err := s.DB.Exec("UPDATE nasabah SET nik = ?, nama = ?, tempat_lahir = ?, tanggal_lahir = ?, alamat = ?, no_telepon = ? WHERE cif = ?",
		nasabah.Nik, nasabah.Nama, nasabah.Tempat_lahir, nasabah.Tanggal_lahir, nasabah.Alamat, nasabah.No_telepon, nasabah.Cif)
	// fmt.Println(nasabah)
	if err != nil {
		panic(err)
	}

	//var Respons *BranchDeliverySystem.NASABAH_INFO
	status, _ := rows.RowsAffected()
	//fmt.Println(idUser)
	if status > 0 {
		return nasabah, nil
	} else {
		return ent.Nasabah{}, nil
	}
}
