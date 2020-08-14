package controller

import (
    "fmt"
    // "html/template"
    // "io"
    "net/http"
    "strconv"

    // ent "bds/entities"
    serv "bds/services"
    conf "bds/configs"

    "github.com/labstack/echo"
)


func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func Login(c echo.Context) error {
    return c.Render(http.StatusOK, "login.html", nil)
}

func LoginProcess(c echo.Context) error {    
    //Koneksi database
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
    id_user, _ := strconv.Atoi(c.FormValue("id_user"))
    password := c.FormValue("password")

    // Melakukan login
    user, err := con.LoginUser(id_user, password)
    if err != nil {
        fmt.Println("Gagal Login")
    }

    // Routing data
    if user.Role == "teller" {
        return c.Render(http.StatusOK, "menu.html", user)
    } else if user.Role == "cs" {
        return c.Render(http.StatusOK, "menu.html", user)
    } else if id_user == 0 || password == "" {
        data := map[string]string{
            "Status" : "Masukan ID User dan Password",
        }
        return c.Render(http.StatusOK, "login.html", data)
    } else {
        data := map[string]string{
            "Status" : "ID User atau Password salah!",
        }
        return c.Render(http.StatusOK, "login.html", data)
    }
}
