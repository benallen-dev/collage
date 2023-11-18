package main

import (
	"io"
	"math/rand"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Template struct {
	templates *template.Template
}

func getRandomName() string {
	names := [7]string {
		"Dave",
		"Peter",
		"Jesus",
		"Mary",
		"Helen",
		"Alice",
		"Bob",
	}

	// make random number
	idx := rand.Intn(7)
	// return that element from the list
	return names[idx]
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
		return c.Render(http.StatusOK, "name", getRandomName())
	})

	e.GET("/foo", func (c echo.Context) error {
		return c.String(http.StatusOK, "foo")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
