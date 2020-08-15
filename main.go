package main

import (
    "html/template"
    "io"
    "net/http"

    // ent "bds/entities"
    cont "bds/controllers"

    "github.com/labstack/echo"
)

// type M map[string]interface{}

type Renderer struct {
    template *template.Template
    debug    bool
    location string
}

func NewRenderer(location string, debug bool) *Renderer {
    tpl := new(Renderer)
    tpl.location = location
    tpl.debug = debug

    tpl.ReloadTemplates()

    return tpl
}

func (t *Renderer) ReloadTemplates() {
    t.template = template.Must(template.ParseGlob(t.location))
}

func (t *Renderer) Render(
    w io.Writer, 
    name string, 
    data interface{}, 
    c echo.Context,
) error {
    if t.debug {
        t.ReloadTemplates()
    }
    return t.template.ExecuteTemplate(w, name, data)
}

func main() {
    e := echo.New()

    e.Renderer = NewRenderer("./views/*.html", true)

    e.Any("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/index")
	})

    e.Any("/index", cont.Index)

    e.GET("/login", cont.Login)

    e.POST("/login_process", cont.LoginProcess)

    e.POST("/setor-tunai", cont.SetorTunai)

    e.Logger.Fatal(e.Start(":9000"))
}

