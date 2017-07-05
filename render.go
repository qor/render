// Package render support to render templates by your control.
package render

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/qor/admin"
)

// DefaultLayout default layout name
const DefaultLayout = "application"

// DefaultViewPath default view path
const DefaultViewPath = "app/views"

// Config render config
type Config struct {
	FuncMapMaker    func(render *Render, request *http.Request, writer http.ResponseWriter) template.FuncMap
	assetFileSystem admin.AssetFSInterface
}

// Render the render struct.
type Render struct {
	*Config

	viewPaths []string
	funcMaps  template.FuncMap
}

// New initalize the render struct.
func New(viewPaths ...string) *Render {
	render := &Render{funcMaps: map[string]interface{}{}, Config: &Config{}}
	render.SetAssetFS(&admin.AssetFileSystem{})

	for _, viewPath := range append(viewPaths, filepath.Join(root, DefaultViewPath)) {
		render.RegisterViewPath(viewPath)
	}

	return render
}

// RegisterViewPath register view path
func (render *Render) RegisterViewPath(paths ...string) {
	for _, pth := range paths {
		if filepath.IsAbs(pth) {
			render.viewPaths = append(render.viewPaths, pth)
			render.assetFileSystem.RegisterPath(pth)
		} else {
			if absPath, err := filepath.Abs(pth); err == nil && isExistingDir(absPath) {
				render.viewPaths = append(render.viewPaths, absPath)
				render.assetFileSystem.RegisterPath(absPath)
			} else if isExistingDir(filepath.Join(root, "vendor", pth)) {
				render.assetFileSystem.RegisterPath(filepath.Join(root, "vendor", pth))
			} else {
				for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
					if p := path.Join(gopath, "src", pth); isExistingDir(p) {
						render.viewPaths = append(render.viewPaths, p)
						render.assetFileSystem.RegisterPath(p)
					}
				}
			}
		}
	}
}

// PrependViewPath prepend view path
func (render *Render) PrependViewPath(paths ...string) {
	for _, pth := range paths {
		if filepath.IsAbs(pth) {
			render.viewPaths = append([]string{pth}, render.viewPaths...)
			render.assetFileSystem.PrependPath(pth)
		} else {
			if absPath, err := filepath.Abs(pth); err == nil && isExistingDir(absPath) {
				render.viewPaths = append([]string{absPath}, render.viewPaths...)
				render.assetFileSystem.PrependPath(absPath)
			} else if isExistingDir(filepath.Join(root, "vendor", pth)) {
				render.assetFileSystem.PrependPath(filepath.Join(root, "vendor", pth))
			} else {
				for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
					if p := path.Join(gopath, "src", pth); isExistingDir(p) {
						render.viewPaths = append([]string{p}, render.viewPaths...)
						render.assetFileSystem.PrependPath(p)
					}
				}
			}
		}
	}
}

// SetAssetFS set asset fs for render
func (render *Render) SetAssetFS(assetFS admin.AssetFSInterface) {
	for _, viewPath := range render.viewPaths {
		assetFS.RegisterPath(viewPath)
	}

	render.assetFileSystem = assetFS
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
