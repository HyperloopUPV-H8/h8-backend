package infra

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestPage(t *testing.T) {
	defaultStaticPath = path.Join("test", "resources")
	t.Cleanup(func() {
		defaultStaticPath = path.Join("static", "build")
	})

	server := New[any, any, any]()
	server.HandleSPA()
	go server.ListenAndServe()

	client := http.Client{}

	type testCase struct {
		path   string
		status int
		body   []byte
	}

	tests := map[string]testCase{
		"index.html": {
			path:   "/",
			status: http.StatusOK,
			body:   readFile(path.Join("test", "resources", "index.html")),
		},
		"index.js": {
			path:   "/index.js",
			status: http.StatusOK,
			body:   readFile(path.Join("test", "resources", "index.js")),
		},
		"invalid": {
			path:   "/foo",
			status: http.StatusNotFound,
			body:   []byte("404 page not found\n"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.Get("http://" + serverAddr + test.path)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != test.status {
				t.Fatalf("expected status %d, got %d", test.status, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(body, test.body) {
				t.Fatalf("expected %s, got %s", test.body, body)
			}
		})
	}
}

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	return data
}
