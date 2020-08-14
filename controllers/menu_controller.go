package controller

import (
    // "html/template"
    // "io"
    "net/http"

    // ent "bds/entities"

    "github.com/labstack/echo"
)


func Menu(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}