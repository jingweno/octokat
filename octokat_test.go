package octokat

import (
	"fmt"
	"github.com/bmizerany/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the GitHub client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a octokat.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// octokat client configured to use test server
	client = NewClient()
	client.BaseURL = server.URL

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m == r.Method && r.URL.String() == "/" {
			respondWith(w, testRootJSON())
		} else {
			http.Error(w, "Bad Request", 400)
		}
	})
}

// teardown closes the test HTTP server.
func tearDown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	assert.Equal(t, want, r.Method)
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	assert.Equal(t, want, r.Header.Get(header))
}

func testBody(t *testing.T, r *http.Request, want string) {
	body, _ := ioutil.ReadAll(r.Body)
	assert.Equal(t, want, string(body))
}

func respondWith(w http.ResponseWriter, s string) {
	fmt.Fprint(w, s)
}

func testURLOf(path string) *url.URL {
	u, _ := url.ParseRequestURI(testURLStringOf(path))
	return u
}

func testURLStringOf(path string) string {
	return fmt.Sprintf("%s/%s", server.URL, path)
}

func testRootJSON() string {
	root := Root{
		CurrentUserURL: Hyperlink(testURLStringOf("user")),
		UserURL:        Hyperlink(testURLStringOf("users/{user}")),
	}
	json, _ := jsonMarshal(root)
	return string(json)
}

func loadFixture(f string) string {
	pwd, _ := os.Getwd()
	p := fmt.Sprintf("%s/fixtures/%s", pwd, f)
	c, _ := ioutil.ReadFile(p)
	return string(c)
}
