# Render

Render Templates (WIP)

## Usage

```go
import "github.com/qor/render"

func main() {
  Render := render.New()
  Render.Layout("application").Render("index", request, write)
}
```

## TODO

* Handle locales
* Response to different content type `index.tmpl`, `index.mobile.tmpl`, `index.mobile+xml.tmpl`
* Bindata
