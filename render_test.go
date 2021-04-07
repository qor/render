package render

import (
	"io/ioutil"
	"regexp"
	"testing"

	"net/http/httptest"
	"net/textproto"
)

func TestExecute(t *testing.T) {
	Render := New(nil, "test")

	request := httptest.NewRequest("GET", "/test", nil)
	responseWriter := httptest.NewRecorder()
	var context interface{}

	tmpl := Render.Layout("layout_for_test")
	tmpl.Execute("test", context, request, responseWriter)

	if textproto.TrimString(responseWriter.Body.String()) != "Template for test" {
		t.Errorf("The template isn't rendered")
	}
}

func TestErrorMessageWhenMissingLayout(t *testing.T) {
	Render := New(nil, "test")

	request := httptest.NewRequest("GET", "/test", nil)
	responseWriter := httptest.NewRecorder()
	var context interface{}

	nonExistLayout := "ThePlant"
	tmpl := Render.Layout(nonExistLayout)
	err := tmpl.Execute("test", context, request, responseWriter)
	if err != nil {
		t.Error("we don't return error, we render the error on page instead")
	}

	bodyBytes, err1 := ioutil.ReadAll(responseWriter.Result().Body)
	if err1 != nil {
		t.Fatal(err1)
	}
	bodyString := string(bodyBytes)

	errorRegexp := "Failed to render page.+" + nonExistLayout + ".*"

	if matched, _ := regexp.MatchString(errorRegexp, bodyString); !matched {
		t.Errorf("Missing layout error message is incorrect")
	}
}

func TestErrorMessageWhenLayoutContainsError(t *testing.T) {
	Render := New(nil, "test")

	request := httptest.NewRequest("GET", "/test", nil)
	responseWriter := httptest.NewRecorder()
	var context interface{}

	layoutContainsError := "layout_contains_error"
	tmpl := Render.Layout(layoutContainsError)
	err := tmpl.Execute("test", context, request, responseWriter)

	if err != nil {
		t.Error("we don't return error, we render the error on page instead")
	}

	bodyBytes, err1 := ioutil.ReadAll(responseWriter.Result().Body)
	if err1 != nil {
		t.Fatal(err1)
	}
	bodyString := string(bodyBytes)

	errorRegexp := "Failed to render page.+" + layoutContainsError + ".*"

	if matched, _ := regexp.MatchString(errorRegexp, bodyString); !matched {
		t.Errorf("Missing layout error message is incorrect")
	}
}
