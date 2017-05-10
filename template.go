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

	if tmpl.layout != "" {
		if b, err := tmpl.findTemplate(filepath.Join("layouts", tmpl.layout)); err == nil {
			if t, err := template.New("").Funcs(funcMap).Parse(string(b)); err == nil {
				return t.Execute(writer, obj)
			} else {
				return err
			}
		} else {
			err := fmt.Errorf("haven't found layout: '%v.tmpl'\n", filepath.Join("layouts", tmpl.layout))
			fmt.Println(err)
			return err
		}
	} else if content, err := tmpl.findTemplate(name); err == nil {
		if t, err := template.New("").Funcs(funcMap).Parse(string(content)); err == nil {
			return t.Execute(writer, obj)
		} else {
			return err
		}
	} else {
		return fmt.Errorf("failed to find template: %v", name)
	}
}

func (tmpl *Template) findTemplate(name string) ([]byte, error) {
	return tmpl.render.AssetFileSystem.Asset(name + ".tmpl")
}
