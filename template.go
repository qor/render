package render

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Template struct {
	render *Render
	layout string
}

func (tmpl *Template) Render(name string, context interface{}, request *http.Request, writer http.ResponseWriter) (err error) {
	if filename, err := tmpl.findTemplate(name); err == nil {
		layoutName, _ := tmpl.findTemplate(filepath.Join("layouts", tmpl.layout))

		if t, err := template.New(filepath.Base(filename)).Funcs(tmpl.render.funcMaps).ParseFiles(layoutName, filename); err == nil {
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
