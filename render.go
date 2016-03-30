package render

import (
	"html/template"
	"net/http"
)

type Render struct {
	ViewPaths []string
	funcMaps  template.FuncMap
}

func New(viewPaths ...string) *Render {
	return &Render{ViewPaths: viewPaths}
}

func (render *Render) Layout(name string) *Template {
	return &Template{render: render, layout: name}
}

func (render *Render) Render(name string, context interface{}, request *http.Request, writer http.ResponseWriter) error {
	return render.Layout("application").Render(name, context, request, writer)
}

func (render *Render) RegisterFuncMap(name string, fc interface{}) {
	render.funcMaps[name] = fc
}
