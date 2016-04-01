package render

import (
	"bytes"
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
	if filename, ok := tmpl.findTemplate(name); ok {
		var filenames = []string{filename}
		if name, ok := tmpl.findTemplate(filepath.Join("layouts", tmpl.layout)); ok {
			filenames = append(filenames, name)
		}

		var t *template.Template
		var funcMap = tmpl.render.funcMaps
		funcMap["render"] = func(name string) template.HTML {
			if filename, ok := tmpl.findTemplate(name); ok {
				result := bytes.NewBufferString("")
				template.New(filepath.Base(filename)).Funcs(funcMap).ParseFiles(filename)
				t.Execute(result, context)
				return template.HTML(result.String())
			}

			fmt.Printf("failed to find template %v\n", name)
			return ""
		}

		if t, err = template.New(filepath.Base(filename)).Funcs(tmpl.render.funcMaps).ParseFiles(filenames...); err == nil {
			return t.Execute(writer, context)
		}
	}

	if err != nil {
		fmt.Printf("Got error when render template %v: %v\n", name, err)
	}
	return err
}

func (tmpl *Template) findTemplate(name string) (string, bool) {
	name = name + ".tmpl"
	for _, viewPath := range tmpl.render.ViewPaths {
		filename := filepath.Join(viewPath, name)
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			return filename, true
		}
	}
	return "", false
}
