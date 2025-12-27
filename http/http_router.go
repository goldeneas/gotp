package http

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type HttpHandler func(request *HTTPRequest, conn net.Conn)

type HttpRouter struct {
	handlers map[string]HttpHandler
}

func NewHttpRouter() *HttpRouter {
	return &HttpRouter{
		handlers: make(map[string]HttpHandler),
	}
}

func (h *HttpRouter) Add(verb string, path string, handler HttpHandler) {
	key := makeKey(verb, path)
	h.handlers[key] = handler
}

func (h *HttpRouter) Call(request *HTTPRequest, conn net.Conn) {
	key := makeKey(request.verb, request.path)
	handler, exists := h.handlers[key]
	if !exists {
		log.Printf("404 Not Found: '%s %s'\n", request.verb, request.path)
		return
	}

	handler(request, conn)
}

func (h *HttpRouter) Remove(verb string, path string) {
	key := makeKey(verb, path)
	delete(h.handlers, key)
}

func (h *HttpRouter) Get(verb string, path string) HttpHandler {
	key := makeKey(verb, path)
	return h.handlers[key]
}

func makeKey(verb string, path string) string {
	return fmt.Sprintf("%s:%s", strings.ToLower(verb), strings.ToLower(path))
}
