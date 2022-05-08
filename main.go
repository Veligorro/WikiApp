package main

import (
	"fmt"
	"github.com/labstack/echo"
	"html/template"
	"io"
	_ "io/ioutil"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

type TemplateRegistry struct {
	templates *template.Template
}

// метод для рендера html шаблонов в echo.
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (p *Page) save() {
	fileName := p.Title + ".txt"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	file.Write(p.Body)
	defer file.Close()
}

func loadPage(title string) *Page {
	fileName := title + ".txt"
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return &Page{Title: title, Body: data}
}

func viewHandler(c echo.Context) error {
	name := c.Param("path")
	p := loadPage(name)
	output := fmt.Sprintf("<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	return c.HTML(http.StatusOK, output)
}

func editHandler(c echo.Context) error {
	name := c.Param("editName")
	p := loadPage(name)

	return c.Render(http.StatusOK, "edit.html", map[string]interface{}{
		"Title": p.Title,
		"Body":  string(p.Body),
	})
}

func main() {
	e := echo.New()
	e.Renderer = &TemplateRegistry{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.GET("/view/:path", viewHandler)
	e.POST("/edit/:editName", editHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
