package render

import (
	"testing"

	"net/http/httptest"
)

func TestExecute(t *testing.T) {
	Render := New("test")

	request := httptest.NewRequest("GET", "/test", nil)
	responseWriter := httptest.NewRecorder()
	var context interface{}

	tmpl := Render.Layout("layout_for_test")
	tmpl.Execute("test", context, request, responseWriter)

	// fmt.Println(tmpl.render.funcMaps["render"]())
	// fmt.Println(responseWriter)

	// t.Errorf("WIP")
}
