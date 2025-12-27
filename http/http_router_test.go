package http

import (
	"net"
	"testing"
)

func TestHttpRouter(t *testing.T) {
	router := NewHttpRouter()
	path := "/api/data"
	verb := "GET"

	called := false
	handler := func(req *HTTPRequest, conn net.Conn) {
		called = true
	}

	router.Add(verb, path, handler)

	// Test Retrieval
	retrieved := router.Get(verb, path)
	if retrieved == nil {
		t.Fatal("Failed to retrieve handler")
	}

	// Test Execution
	request := &HTTPRequest{
		verb:    verb,
		path:    path,
		content: "",
		headers: nil,
	}

	router.Call(request, nil)
	if !called {
		t.Error("Handler was not executed")
	}

	// Test Removal
	router.Remove(verb, path)
	if router.Get(verb, path) != nil {
		t.Error("Handler should have been removed")
	}
}
