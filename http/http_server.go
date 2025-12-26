package http

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
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
	address := conn.RemoteAddr().String()

	for {
		deadline := time.Now().Add(5 * time.Second)
		conn.SetReadDeadline(deadline)

		lines, err := readLines(reader)
		if !handleError(err, conn) {
			break
		}

		// this can happen if the client closes the connection
		// while waiting for an answer (e.g. timeout reached)
		// a client sends '\r\n' which triggers our parsing to stop immediately
		if len(lines) == 0 {
			break
		}

		verb := extractVerb(lines)
		path := extractPath(lines)
		headers := extractHeaders(lines)

		content, err := readContent(headers, reader)
		if !handleError(err, conn) {
			break
		}

		h.router.Call(verb, path, headers, content, conn)
	}

	log.Printf("closing connection with %s", address)
}

func handleError(err error, conn net.Conn) bool {
	if err == nil {
		return true
	}

	address := conn.RemoteAddr().String()

	// we hit a timeout, disconnect
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		log.Printf("hit a timeout for %s: disconnecting", address)
		return false
	}

	// the client wants to disconnect
	if err == io.EOF {
		log.Printf("client at %s sent eof: disconnecting", address)
		return false
	}

	// general error, disconnect
	log.Printf("error while handling communication: %s", err)
	return false
}
