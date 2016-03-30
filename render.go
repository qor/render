package render

import "net/http"

type Render struct {
	ViewPaths string
}

func New(viewPaths ...string) *Render {
	return &Render{ViewPaths: viewPaths}
}

type Template struct {
	*Render
	Layout string
}

func (render *Render) Layout(name string) *Template {
	return Template{Render: render, Layout: name}
}

func (render *Render) Render(name string) *Template {
	return render.Layout("application").Render(name)
}

func (tmpl *Template) Render(name string, context interface{}, request *http.Request, writer http.ResponseWriter) *Template {
	return tmpl.Execute(writer, context)
}
