package entities

type Pesan struct {
	Nama string
}

type User struct {
	Id_user 	int64
	Password 	string
	Nama_user 	string
	Role		string
	Cabang		string
}


type Nasabah struct {
	Cif     		int64
	Nik     		int64
	Nama    		string
	Tempat_lahir	string
	Tanggal_lahir 	string
	Alamat  		string
	No_telepon 		string
}

type NasabahDetail struct {
	Cif			int64
	Nama		string
	No_rekening	int64
	Saldo		int64
}

type Transaksi struct{
	Id_transaksi	int64
	Id_user			int64
	No_rekening		int64
	Tanggal			string
	Jenis_transaksi	string
	Nominal			int64
	Saldo			int64
	Berita			string
}


type TransaksiPB struct {
	Id_user 		int64
	Tanggal			string
	Nasabah_debit	NasabahDetail
	Nasabah_kredit	NasabahDetail
	Nominal			int64
	Berita			string
}