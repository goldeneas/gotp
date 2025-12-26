package http

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type config struct {
	logs bool
}

type HttpServer struct {
	port     int
	config   config
	listener net.Listener
	router   *HttpRouter
}

type ConfigFn func(c *config)

func SetEnableLogs(enable bool) ConfigFn {
	return func(c *config) {
		c.logs = enable
	}
}

func NewServer(router *HttpRouter, opts ...ConfigFn) *HttpServer {
	config := config{
		logs: false,
	}

	for _, f := range opts {
		f(&config)
	}

	return &HttpServer{
		config: config,
		router: router,
	}
}

func (h *HttpServer) Listen(port int) error {
	if h.listener != nil {
		h.listener.Close()
	}

	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	h.listener = listener

	if h.config.logs {
		log.Printf("server listening on address '%s'", address)
	}

	return nil
}

func (h *HttpServer) Accept() error {
	conn, err := h.listener.Accept()
	if err != nil {
		return err
	}

	if h.config.logs {
		log.Printf("server accepted connection from: '%s'", conn.RemoteAddr().String())
	}

	go h.connectionHandler(conn)
	return nil
}

func (h *HttpServer) Close() error {
	return h.listener.Close()
}

func (h *HttpServer) connectionHandler(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	lines := extractLines(reader)

	verb := extractVerb(lines)
	path := extractPath(lines)
	headers := extractHeaders(lines)

	content, err := readContent(headers, reader)
	if err != nil {
		log.Printf("error while extracting content: %s", err)
		return
	}

	h.router.Call(verb, path, headers, content)
}
