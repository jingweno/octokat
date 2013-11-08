package octokit

import (
	"github.com/bmizerany/assert"
	"net/http"
	"strings"
	"testing"
)

func TestResponseError_Error_400(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		respondWith(w, `{"message":"Problems parsing JSON"}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Tf(t, strings.Contains(err.Error(), "400 - Problems parsing JSON"), "%s", err.Error())

	e := err.(*ResponseError)
	assert.Equal(t, ErrorBadRequest, e.Type)
}

func TestResponseError_Error_422_error(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(422)
		respondWith(w, `{"error":"No repository found for hubtopic"}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Tf(t, strings.Contains(err.Error(), "Error: No repository found for hubtopic"), "%s", err.Error())

	e := err.(*ResponseError)
	assert.Equal(t, ErrorUnprocessableEntity, e.Type)
}

func TestResponseError_Error_422_error_summary(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(422)
		respondWith(w, `{"message":"Validation Failed", "errors": [{"resource":"Issue", "field": "title", "code": "missing_field"}]}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Tf(t, strings.Contains(err.Error(), "422 - Validation Failed"), "%s", err.Error())
	assert.Tf(t, strings.Contains(err.Error(), "missing_field error caused by title field on Issue resource"), "%s", err.Error())

	e := err.(*ResponseError)
	assert.Equal(t, ErrorUnprocessableEntity, e.Type)
}

func TestResponseError_Error_415(t *testing.T) {
	setup()
	defer tearDown()

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		respondWith(w, `{"message":"Unsupported Media Type", "documentation_url":"http://developer.github.com/v3"}`)
	})

	req, _ := client.NewRequest("error")
	_, err := req.Get(nil)
	assert.Tf(t, strings.Contains(err.Error(), "415 - Unsupported Media Type"), "%s", err.Error())
	assert.Tf(t, strings.Contains(err.Error(), "// See: http://developer.github.com/v3"), "%s", err.Error())

	e := err.(*ResponseError)
	assert.Equal(t, ErrorUnsupportedMediaType, e.Type)
}
