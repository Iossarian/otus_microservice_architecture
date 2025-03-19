package build

import (
	"context"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (b *Builder) RestServer() *echo.Echo {
	e := echo.New()

	b.shutdown.add(func(ctx context.Context) error {
		return e.Close()
	})

	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler, err := b.handler()
	if err != nil {
		e.Logger.Fatal(err)
	}

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = renderer

	e.GET("/orders", handler.OrderForm)
	e.POST("/orders", handler.Create)
	e.GET("/orders/:id", handler.Get)
	e.GET("/key/:key", handler.DeleteKey)

	return e
}
