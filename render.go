package render

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// Render render struct
type Render struct {
	ViewPaths []string
	funcMaps  template.FuncMap
}

// New initalize render struct
func New(viewPaths ...string) *Render {
	if isExistingDir(filepath.Join(root, "app/views")) {
		viewPaths = append(viewPaths, filepath.Join(root, "app/views"))
	}

	return &Render{ViewPaths: viewPaths, funcMaps: map[string]interface{}{}}
}

// Layout set layout for tmpl
func (render *Render) Layout(name string) *Template {
	return &Template{render: render, layout: name}
}

// Funcs set funcs for tmpl with application layout
func (render *Render) Funcs(funcMap template.FuncMap) *Template {
	return render.Layout("application").Funcs(funcMap)
}

// Execute execute tmpl with application layout
func (render *Render) Execute(name string, context interface{}, request *http.Request, writer http.ResponseWriter) error {
	return render.Layout("application").Execute(name, context, request, writer)
}

// RegisterFuncMap register FuncMap for render
func (render *Render) RegisterFuncMap(name string, fc interface{}) {
	if render.funcMaps == nil {
		render.funcMaps = template.FuncMap{}
	}
	render.funcMaps[name] = fc
}
