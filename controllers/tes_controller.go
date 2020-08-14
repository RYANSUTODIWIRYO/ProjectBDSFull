package controller

import (
    // "html/template"
    // "io"
    "net/http"

    ent "bds/entities"

    "github.com/labstack/echo"
)


func Index(c echo.Context) error {
	data := ent.Pesan{
		Nama : "Ryan",
	}
	return c.Render(http.StatusOK, "index.html", data)
}