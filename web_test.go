package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoHandler(t *testing.T) {
	for _, test := range []struct {
		url    string
		header http.Header
		status int
		body   string
	}{
		{
			url:    "/",
			status: http.StatusOK,
			body:   "URL: /\nHeader:\n",
		},
		{
			url:    "/",
			header: http.Header{},
			status: http.StatusOK,
			body:   "URL: /\nHeader:\n",
		},
		{
			url: "/",
			header: http.Header{
				"key1": []string{"value1"},
			},
			status: http.StatusOK,
			body:   "URL: /\nHeader:\nkey1 -> \"value1\"\n",
		},
		{
			url: "/",
			header: http.Header{
				"key1": []string{"value1", "value2"},
			},
			status: http.StatusOK,
			body:   "URL: /\nHeader:\nkey1 -> \"value1\"; \"value2\"\n",
		},
		{
			url: "/",
			header: http.Header{
				"key1": []string{"value1"},
				"key2": []string{"value2"},
			},
			status: http.StatusOK,
			body:   "URL: /\nHeader:\nkey1 -> \"value1\"\nkey2 -> \"value2\"\n",
		},
		{
			url: "/",
			header: http.Header{
				"key2": []string{"value2"},
				"key1": []string{"value1"},
			},
			status: http.StatusOK,
			body:   "URL: /\nHeader:\nkey1 -> \"value1\"\nkey2 -> \"value2\"\n",
		},
	} {
		req, err := http.NewRequest(http.MethodGet, test.url, nil)
		if err != nil {
			t.Fatalf("can not create request: %s", err)
		}
		req.Header = test.header

		w := httptest.NewRecorder()

		echoHandler(w, req)

		if w.Code != test.status {
			t.Errorf("got status %d, expected %d", w.Code, test.status)
		}

		body := w.Body.String()
		if body != test.body {
			t.Errorf("got '%s', wanted '%s'", body, test.body)
		}
	}
}

func TestVersionHandler(t *testing.T) {
	expectedStatus := http.StatusOK
	expectedContent := "application/json; chatset=utf-8"

	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("can not create request: %s", err)
	}

	w := httptest.NewRecorder()
	handler := versionHandler("test-version")

	handler(w, req)

	if w.Code != expectedStatus {
		t.Errorf("got status %d, wanted %d", w.Code, expectedStatus)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != expectedContent {
		t.Errorf("got type '%s', wanted '%s'", contentType, expectedContent)
	}

	body := w.Body.String()
	if len(body) == 0 {
		t.Error("got no body.")
	}
}

func TestCreateServer(t *testing.T) {
	addr := "localhost:3000"
	s := createServer(addr, "")

	if s.Addr != addr {
		t.Errorf("got '%s', wanted '%s'", s.Addr, addr)
	}

	if s.Handler == nil {
		t.Errorf("handler not set")
	}
}
