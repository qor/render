package render

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Render struct {
	ViewPaths []string
}

func New(viewPaths ...string) *Render {
	return &Render{ViewPaths: viewPaths}
}

type Template struct {
	render *Render
	layout string
}

func (render *Render) Layout(name string) *Template {
	return &Template{render: render, layout: name}
}

func (render *Render) Render(name string, context interface{}, request *http.Request, writer http.ResponseWriter) error {
	return render.Layout("application").Render(name, context, request, writer)
}

func (tmpl *Template) Render(name string, context interface{}, request *http.Request, writer http.ResponseWriter) (err error) {
	if filename, err := tmpl.findTemplate(name); err == nil {
		layoutName, _ := tmpl.findTemplate(filepath.Join("layouts", tmpl.layout))

		if t, err := template.New(filepath.Base(filename)).ParseFiles(layoutName, filename); err == nil {
			return t.Execute(writer, context)
		}
	}
	return err
}

func (tmpl *Template) findTemplate(name string) (string, error) {
	for _, viewPath := range tmpl.render.ViewPaths {
		filename := filepath.Join(viewPath, name)
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			return filename, nil
		}
	}
	return "", fmt.Errorf("template not found: %v", name)
}
