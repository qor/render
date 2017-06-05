// Package render support to render templates by your control.
package render

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/qor/admin"
)

// DefaultLayout default layout name
const DefaultLayout = "application"

// DefaultViewPath default view path
const DefaultViewPath = "app/views"

// Render the render struct.
type Render struct {
	AssetFileSystem admin.AssetFSInterface

	viewPaths []string
	funcMaps  template.FuncMap
}

// New initalize the render struct.
func New(viewPaths ...string) *Render {
	if isExistingDir(filepath.Join(root, DefaultViewPath)) {
		viewPaths = append(viewPaths, filepath.Join(root, DefaultViewPath))
	}

	render := &Render{viewPaths: viewPaths, funcMaps: map[string]interface{}{}}
	render.SetAssetFS(&admin.AssetFileSystem{})

	return render
}

// RegisterViewPath register view path
func (render *Render) RegisterViewPath(pth string) {
	render.viewPaths = append(render.viewPaths, pth)
	render.AssetFileSystem.RegisterPath(pth)
}

// SetAssetFS set asset fs for render
func (render *Render) SetAssetFS(assetFS admin.AssetFSInterface) {
	for _, viewPath := range render.viewPaths {
		assetFS.RegisterPath(viewPath)
	}

	render.AssetFileSystem = assetFS
}

// Layout set layout for template.
func (render *Render) Layout(name string) *Template {
	return &Template{render: render, layout: name}
}

// Funcs set helper functions for template with default "application" layout.
func (render *Render) Funcs(funcMap template.FuncMap) *Template {
	return render.Layout(DefaultLayout).Funcs(funcMap)
}

// Execute render template with default "application" layout.
func (render *Render) Execute(name string, context interface{}, request *http.Request, writer http.ResponseWriter) error {
	return render.Layout(DefaultLayout).Execute(name, context, request, writer)
}

// RegisterFuncMap register FuncMap for render.
func (render *Render) RegisterFuncMap(name string, fc interface{}) {
	if render.funcMaps == nil {
		render.funcMaps = template.FuncMap{}
	}
	render.funcMaps[name] = fc
}
