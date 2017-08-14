// Package render support to render templates by your control.
package render

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/qor/assetfs"
	"github.com/qor/qor/utils"
)

// DefaultLayout default layout name
const DefaultLayout = "application"

// DefaultViewPath default view path
const DefaultViewPath = "app/views"

// Config render config
type Config struct {
	IgnoreLayoutError bool
	ViewPaths         []string
	DefaultLayout     string
	FuncMapMaker      func(render *Render, request *http.Request, writer http.ResponseWriter) template.FuncMap
	assetFileSystem   assetfs.Interface
}

// Render the render struct.
type Render struct {
	*Config

	funcMaps template.FuncMap
}

// New initalize the render struct.
func New(config *Config, viewPaths ...string) *Render {
	if config == nil {
		config = &Config{}
	}

	if config.DefaultLayout != "" {
		config.DefaultLayout = DefaultLayout
	}

	render := &Render{funcMaps: map[string]interface{}{}, Config: config}
	render.SetAssetFS(assetfs.AssetFS().NameSpace("views"))

	config.ViewPaths = append(append(config.ViewPaths, viewPaths...), DefaultViewPath)

	for _, viewPath := range config.ViewPaths {
		render.RegisterViewPath(viewPath)
	}

	return render
}

// RegisterViewPath register view path
func (render *Render) RegisterViewPath(paths ...string) {
	for _, pth := range paths {
		if filepath.IsAbs(pth) {
			render.ViewPaths = append(render.ViewPaths, pth)
			render.assetFileSystem.RegisterPath(pth)
		} else {
			if absPath, err := filepath.Abs(pth); err == nil && isExistingDir(absPath) {
				render.ViewPaths = append(render.ViewPaths, absPath)
				render.assetFileSystem.RegisterPath(absPath)
			} else if isExistingDir(filepath.Join(utils.AppRoot, "vendor", pth)) {
				render.assetFileSystem.RegisterPath(filepath.Join(utils.AppRoot, "vendor", pth))
			} else {
				for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
					if p := path.Join(gopath, "src", pth); isExistingDir(p) {
						render.ViewPaths = append(render.ViewPaths, p)
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
			render.ViewPaths = append([]string{pth}, render.ViewPaths...)
			render.assetFileSystem.PrependPath(pth)
		} else {
			if absPath, err := filepath.Abs(pth); err == nil && isExistingDir(absPath) {
				render.ViewPaths = append([]string{absPath}, render.ViewPaths...)
				render.assetFileSystem.PrependPath(absPath)
			} else if isExistingDir(filepath.Join(utils.AppRoot, "vendor", pth)) {
				render.assetFileSystem.PrependPath(filepath.Join(utils.AppRoot, "vendor", pth))
			} else {
				for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
					if p := path.Join(gopath, "src", pth); isExistingDir(p) {
						render.ViewPaths = append([]string{p}, render.ViewPaths...)
						render.assetFileSystem.PrependPath(p)
					}
				}
			}
		}
	}
}

// SetAssetFS set asset fs for render
func (render *Render) SetAssetFS(assetFS assetfs.Interface) {
	for _, viewPath := range render.ViewPaths {
		assetFS.RegisterPath(viewPath)
	}

	render.assetFileSystem = assetFS
}

// Layout set layout for template.
func (render *Render) Layout(name string) *Template {
	return &Template{render: render, layout: name, ignoreLayoutError: render.IgnoreLayoutError}
}

// Funcs set helper functions for template with default "application" layout.
func (render *Render) Funcs(funcMap template.FuncMap) *Template {
	return render.Layout(render.Config.DefaultLayout).Funcs(funcMap)
}

// Execute render template with default "application" layout.
func (render *Render) Execute(name string, context interface{}, request *http.Request, writer http.ResponseWriter) error {
	return render.Layout(render.Config.DefaultLayout).Execute(name, context, request, writer)
}

// RegisterFuncMap register FuncMap for render.
func (render *Render) RegisterFuncMap(name string, fc interface{}) {
	if render.funcMaps == nil {
		render.funcMaps = template.FuncMap{}
	}
	render.funcMaps[name] = fc
}

// Asset get content from AssetFS by name
func (render *Render) Asset(name string) ([]byte, error) {
	return render.assetFileSystem.Asset(name)
}
