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
	render  *Render
	layout  string
	funcMap template.FuncMap
}

// FuncMap get func maps from tmpl
func (tmpl *Template) FuncMap() template.FuncMap {
	if tmpl.funcMap == nil {
		return tmpl.render.funcMaps
	}

	var funcMap = tmpl.funcMap
	for key, value := range tmpl.render.funcMaps {
		funcMap[key] = value
	}
	return funcMap
}

// Funcs register Funcs for tmpl
func (tmpl *Template) Funcs(funcMap template.FuncMap) *Template {
	tmpl.funcMap = funcMap
	return tmpl
}

// Execute execute tmpl
func (tmpl *Template) Execute(name string, context interface{}, request *http.Request, writer http.ResponseWriter) (err error) {
	if filename, ok := tmpl.findTemplate(name); ok {
		// filenames
		var filenames []string
		var layout string
		if layout, ok = tmpl.findTemplate(filepath.Join("layouts", tmpl.layout)); ok {
			filenames = append(filenames, layout)
		}
		// append templates to last, then it could be used to overwrite layouts templates
		filenames = append(filenames, filename)

		var obj = map[string]interface{}{
			"Template": name,
			"Result":   context,
		}

		// funcMaps
		var funcMap = tmpl.FuncMap()
		funcMap["render"] = func(name string, objs ...interface{}) (template.HTML, error) {
			var (
				err       error
				renderObj interface{}
			)

			if len(objs) == 0 {
				// default obj
				renderObj = obj
			} else {
				// overwrite obj
				for _, o := range objs {
					renderObj = o
					break
				}
			}

			if filename, ok := tmpl.findTemplate(name); ok {
				var partialTemplate *template.Template
				result := bytes.NewBufferString("")
				if partialTemplate, err = template.New(filepath.Base(filename)).Funcs(funcMap).ParseFiles(filename); err == nil {
					if err = partialTemplate.Execute(result, renderObj); err == nil {
						return template.HTML(result.String()), err
					}
				}
			} else {
				err = fmt.Errorf("failed to find template: %v", name)
			}

			return "", err
		}

		// parse templates
		var t *template.Template
		if t, err = template.New(filepath.Base(layout)).Funcs(funcMap).ParseFiles(filenames...); err == nil {
			err = t.Execute(writer, obj)
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
