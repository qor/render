# Render

Render Templates

## Usage

```go
import "github.com/qor/render"

func main() {
  Render := render.New()

  Render.Execute("index", obj, request, writer)

  Render.Layout("application").Execute("index", obj, request, writer)

  Render.Layout("application").Funcs(funcsMap).Execute("index", obj, request, writer)
}
```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).
