package http

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const KEEPALIVE_TIMEOUT = 5 * time.Second

type config struct {
	logs         bool
	serveEnabled bool
	servePath    string
}

type HttpServer struct {
	port       int
	config     config
	listener   net.Listener
	router     *HttpRouter
	fileServer *HttpFileServer
}

type ConfigFn func(c *config)

func SetEnableLogs(enable bool) ConfigFn {
	return func(c *config) {
		c.logs = enable
	}
}

func SetServeEnabled(enable bool) ConfigFn {
	return func(c *config) {
		c.serveEnabled = enable
	}
}

func SetServePath(path string) ConfigFn {
	return func(c *config) {
		c.servePath = path
	}
}

func NewServer(router *HttpRouter, opts ...ConfigFn) *HttpServer {
	config := config{
		logs:         false,
		serveEnabled: false,
		servePath:    "static/",
	}

	for _, f := range opts {
		f(&config)
	}

	return &HttpServer{
		config:     config,
		router:     router,
		fileServer: NewHttpFileServer(),
	}
}

func (h *HttpServer) Listen(port int) error {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	h.listener = listener

	servePath := h.config.servePath
	go h.fileServer.WatchDir(servePath)

	if h.config.logs {
		log.Printf("server listening on address '%s'", address)
		log.Printf("server watching dir '%s'", servePath)
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
		deadline := time.Now().Add(KEEPALIVE_TIMEOUT)
		conn.SetReadDeadline(deadline)

		request, err := h.readRequest(address, reader)
		if err != nil {
			logRequestError(address, err)
			break
		}

		file, err := h.readFile(request.Path())
		if err != nil {
			logFileError(address, err)
		}

		if file != nil {
			Serve(file, "text/html", conn)
		} else {
			h.router.Call(request, conn)
		}

		if request.IsConnectionClose() {
			break
		}
	}

	if h.config.logs {
		log.Printf("closing connection with %s", address)
	}
}

func (h *HttpServer) readFile(path string) ([]byte, error) {
	basePath := h.config.servePath

	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(absBase, path)

	absTarget, err := filepath.Abs(fullPath)
	if err != nil {
		return nil, err
	}

	prefix := absBase + string(filepath.Separator)
	if !strings.HasPrefix(absTarget, prefix) && absTarget != absBase {
		return nil, fmt.Errorf("403 Forbidden: illegal path traversal attempt")
	}

	return h.fileServer.Get(absTarget)
}

func (h *HttpServer) readRequest(address string, reader *bufio.Reader) (*HTTPRequest, error) {
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
		logContentError(address, err)
	}

	path, queries := extractPathAndQueries(lines)

	return &HTTPRequest{
		verb:    extractVerb(lines),
		path:    path,
		queries: queries,
		headers: headers,
		content: content,
	}, nil
}

func logRequestError(address string, err error) {
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

func logFileError(address string, err error) {
	if errors.Is(err, os.ErrNotExist) {
		return
	}

	log.Printf("error while handling static file serving with %s: %s", address, err)
}

func logContentError(address string, err error) {
	if errors.Is(err, strconv.ErrSyntax) {
		return
	}

	log.Printf("error while handling incoming content from %s: %s", address, err)
}
