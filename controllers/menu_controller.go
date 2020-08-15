package controller

import (
	conf "bds/configs"
	ent "bds/entities"
	serv "bds/services"
	"time"

	// "html/template"
	// "io"
	"net/http"
	"strconv"

	// ent "bds/entities"

	"github.com/labstack/echo"
)

func SetorTunai(c echo.Context) error {
	// fmt.Println("Proses Login..")
	// Koneksi database
	db, err := conf.KoneksiDB()
	if err != nil {
		panic(err)
		// return ent.User{}, err
	}
	defer db.Close()

	// Membuat struct koneksi
	con := serv.Service{
		db,
	}

	// Mengambil nilai dari form html
	rek_tujuan, _ := strconv.Atoi(c.FormValue("rek_tujuan"))
	nominal, _ := strconv.Atoi(c.FormValue("nominal"))
	berita := c.FormValue("berita")
	tanggal := time.Now().Format("2006-01-02 15:04:05")

	transaksi := ent.Transaksi{
		No_rekening: int64(rek_tujuan),
		Berita:      berita,
		Nominal:     int64(nominal),
		Tanggal:     tanggal,
	}
	// setor tunai
	status, data, err := con.SetorTunaiService(transaksi)
	if err != nil {
		panic(err)
	}
	if status > 0 {
		return c.Render(http.StatusOK, "teller.html", data)
	} else {
		data := map[string]string{
			"Status": "Setor tunai gagal ",
		}
		return c.Render(http.StatusOK, "menu.html", data)
	}

}

func Teller(c echo.Context) error {
    return c.Render(http.StatusOK, "teller.html", nil)
}