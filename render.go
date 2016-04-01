package render

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Render struct {
	ViewPaths []string
	funcMaps  template.FuncMap
}

func New(viewPaths ...string) *Render {
	if isExistingDir(filepath.Join(root, "app/views")) {
		viewPaths = append(viewPaths, filepath.Join(root, "app/views"))
	}

	return &Render{ViewPaths: viewPaths}
}

func (render *Render) Layout(name string) *Template {
	return &Template{render: render, layout: name}
}

func (render *Render) Render(name string, context interface{}, request *http.Request, writer http.ResponseWriter) error {
	return render.Layout("application").Render(name, context, request, writer)
}

func (render *Render) RegisterFuncMap(name string, fc interface{}) {
	if render.funcMaps == nil {
		render.funcMaps = template.FuncMap{}
	}
	render.funcMaps[name] = fc
}
