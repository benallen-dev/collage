package main

import (
	"io"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/benallen-dev/collage/pkg/util"
)

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Template struct {
	templates *template.Template
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.Static("/", "static")

	e.PUT("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/name", func(c echo.Context) error {
		return c.Render(http.StatusOK, "name", util.GetRandomName())
	})

	e.GET("/foo", func (c echo.Context) error {
		return c.String(http.StatusOK, "foo")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
