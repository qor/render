package render

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// Template template struct
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
	if content, err := tmpl.findTemplate(name); err == nil {
		// filenames
		var (
			contents      []string
			layoutContent string = fmt.Sprintf("{{render %q}}", name)
		)

		layoutPath := filepath.Join("layouts", tmpl.layout)

		if b, err := tmpl.findTemplate(layoutPath); err == nil {
			layoutContent = string(b)
		} else {
			if absoluteLayoutPath, pathErr := filepath.Abs(layoutPath); pathErr == nil {
				err = fmt.Errorf("Cannot find layout: '%v.tmpl'", absoluteLayoutPath)
			} else {
				err = fmt.Errorf("Cannot find layout: '%v.tmpl'", layoutPath)
			}

			fmt.Println("Got error when finding layout:", err)
			return err
		}

		// append templates to last, then it could be used to overwrite layouts templates
		contents = append(contents, string(content))

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

			if renderContent, err := tmpl.findTemplate(name); err == nil {
				var partialTemplate *template.Template
				result := bytes.NewBufferString("")
				if partialTemplate, err = template.New(filepath.Base(name)).Funcs(funcMap).Parse(string(renderContent)); err == nil {
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
		if t, err = template.New("").Funcs(funcMap).Parse(layoutContent); err == nil {
			err = t.Execute(writer, obj)
		}
	} else {
		err = fmt.Errorf("failed to find template: %v", name)
	}

	if err != nil {
		fmt.Printf("Got error when render template %v: %v\n", name, err)
	}
	return err
}

func (tmpl *Template) findTemplate(name string) ([]byte, error) {
	return tmpl.render.AssetFileSystem.Asset(name + ".tmpl")
}
