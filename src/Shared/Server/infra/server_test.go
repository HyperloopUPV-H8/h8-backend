package infra

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/server/infra/interfaces"
	"github.com/gorilla/websocket"
)

func TestPage(t *testing.T) {
	defaultStaticPath = path.Join("test", "resources")
	serverAddr = "127.0.0.1:4000"

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
			defer resp.Body.Close()

			if !reflect.DeepEqual(body, test.body) {
				t.Fatalf("expected %s, got %s", test.body, body)
			}
		})
	}
}

func TestLog(t *testing.T) {
	defaultStaticPath = path.Join("test", "resources")
	serverAddr = "127.0.0.1:4001"

	logChan := make(chan bool)
	server := New[any, any, any]()
	server.HandleLog("/backend/log", logChan)
	go server.ListenAndServe()

	client := http.Client{}

	resp := put(client, "/backend/log", "enable")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	select {
	case enable := <-logChan:
		if enable != true {
			log.Fatalln("expected enable")
		}
	case <-time.After(time.Millisecond):
		log.Fatalln("expected response")
	}

	resp = put(client, "/backend/log", "disable")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	select {
	case enable := <-logChan:
		if enable != false {
			log.Fatalln("expected disable")
		}
	case <-time.After(time.Millisecond):
		log.Fatalln("expected response")
	}

	resp = put(client, "/backend/log", "foo")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	select {
	case <-logChan:
		log.Fatalln("expected nothing")
	case <-time.After(time.Millisecond):
	}
}

func TestWebSocket(t *testing.T) {
	defaultStaticPath = path.Join("test", "resources")
	serverAddr = "127.0.0.1:4002"

	server := New[any, any, any]()
	server.HandleWebSocketData("/backend/data", func(interfaces.WebSocket, <-chan any) {})
	server.HandleWebSocketMessage("/backend/message", func(interfaces.WebSocket, <-chan any) {})
	server.HandleWebSocketOrder("/backend/order", func(interfaces.WebSocket, chan<- any) {})
	go server.ListenAndServe()

	_, resp := wsClient("ws://" + serverAddr + "/backend/data")
	resp.Body.Close()

	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("expected status %d, got %d", http.StatusSwitchingProtocols, resp.StatusCode)
	}

	_, resp = wsClient("ws://" + serverAddr + "/backend/message")
	resp.Body.Close()

	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("expected status %d, got %d", http.StatusSwitchingProtocols, resp.StatusCode)
	}

	_, resp = wsClient("ws://" + serverAddr + "/backend/order")
	resp.Body.Close()

	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("expected status %d, got %d", http.StatusSwitchingProtocols, resp.StatusCode)
	}
}

func put(client http.Client, path string, body string) *http.Response {
	resp, err := client.Do(&http.Request{
		Method:        http.MethodPut,
		Body:          io.NopCloser(bytes.NewBuffer([]byte(body))),
		ContentLength: int64(len(body)),
		URL:           getURL("http://" + serverAddr + path),
		Proto:         "HTTP/1.1",
		ProtoMinor:    1,
		ProtoMajor:    1,
	})

	if err != nil {
		log.Fatalln(err)
	}

	return resp
}

func wsClient(path string) (*websocket.Conn, *http.Response) {
	conn, resp, err := websocket.DefaultDialer.Dial(path, http.Header{})
	if err != nil {
		log.Fatalln(err)
	}

	return conn, resp
}

func getURL(path string) *url.URL {
	result, err := url.Parse(path)
	if err != nil {
		log.Fatalln(err)
	}

	return result
}

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	return data
}
