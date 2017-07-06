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
func (tmpl *Template) funcMapMaker(req *http.Request, writer http.ResponseWriter) template.FuncMap {
	var funcMap = template.FuncMap{}

	for key, fc := range tmpl.render.funcMaps {
		funcMap[key] = fc
	}

	if tmpl.render.Config.FuncMapMaker != nil {
		for key, fc := range tmpl.render.Config.FuncMapMaker(tmpl.render, req, writer) {
			funcMap[key] = fc
		}
	}

	for key, fc := range tmpl.funcMap {
		funcMap[key] = fc
	}
	return funcMap
}

// Funcs register Funcs for tmpl
func (tmpl *Template) Funcs(funcMap template.FuncMap) *Template {
	tmpl.funcMap = funcMap
	return tmpl
}

// Execute execute tmpl
func (tmpl *Template) Execute(name string, context interface{}, request *http.Request, writer http.ResponseWriter) error {
	var obj = map[string]interface{}{
		"Template": name,
		"Result":   context,
	}

	// funcMaps
	var funcMap = tmpl.funcMapMaker(request, writer)
	funcMap["render"] = func(name string, objs ...interface{}) (template.HTML, error) {
		var (
			err           error
			renderObj     interface{}
			renderContent []byte
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

		if renderContent, err = tmpl.findTemplate(name); err == nil {
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

		if err != nil {
			fmt.Println(err)
		}

		return "", err
	}

	var (
		content []byte
		t       *template.Template
		err     error
	)

	if tmpl.layout != "" {
		if content, err = tmpl.findTemplate(filepath.Join("layouts", tmpl.layout)); err == nil {
			if t, err = template.New("").Funcs(funcMap).Parse(string(content)); err == nil {
				var tpl bytes.Buffer
				if err = t.Execute(&tpl, obj); err == nil {
					_, err = writer.Write(tpl.Bytes())
				}
			}
		} else {
			err = fmt.Errorf("haven't found layout: '%v.tmpl'", filepath.Join("layouts", tmpl.layout))
		}
	} else if content, err = tmpl.findTemplate(name); err == nil {
		if t, err = template.New("").Funcs(funcMap).Parse(string(content)); err == nil {
			var tpl bytes.Buffer
			if err = t.Execute(&tpl, obj); err == nil {
				_, err = writer.Write(tpl.Bytes())
			}
		}
	} else {
		err = fmt.Errorf("failed to find template: %v", name)
	}

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (tmpl *Template) findTemplate(name string) ([]byte, error) {
	return tmpl.render.assetFileSystem.Asset(name + ".tmpl")
}
