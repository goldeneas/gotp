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

		request, err := h.readRequest(reader)
		if err != nil {
			logError(address, err)
			break
		}

		h.router.Call(request, conn)

		if request.IsConnectionClose() {
			break
		}
	}

	if h.config.logs {
		log.Printf("closing connection with %s", address)
	}
}

func (h *HttpServer) readRequest(reader *bufio.Reader) (*HttpRequest, error) {
	lines, err := readLines(reader)
	if err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, io.EOF
	}

	headers := extractHeaders(lines)
	content, err := readContent(headers, reader)

	if err != nil {
		return nil, err
	}

	return &HttpRequest{
		verb:    extractVerb(lines),
		path:    extractPath(lines),
		headers: headers,
		content: content,
	}, nil
}

func logError(address string, err error) {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		log.Printf("hit a timeout for %s: disconnecting", address)
		return
	}

	if errors.Is(err, io.EOF) {
		log.Printf("client at %s sent eof: disconnecting", address)
		return
	}

	log.Printf("error while handling communication with %s: %s", address, err)
}
