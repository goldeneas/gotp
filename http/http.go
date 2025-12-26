package http

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type config struct {
	logs bool
}

type HttpServer struct {
	port     int
	config   config
	listener net.Listener
}

type HttpMessageHandler func(verb string, content string, headers map[string]string)

type ConfigFn func(c *config)

func SetEnableLogs(enable bool) ConfigFn {
	return func(c *config) {
		c.logs = enable
	}
}

func NewServer(opts ...ConfigFn) *HttpServer {
	config := config{
		logs: false,
	}

	for _, f := range opts {
		f(&config)
	}

	return &HttpServer{
		config: config,
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

func (h *HttpServer) Accept(handler HttpMessageHandler) error {
	conn, err := h.listener.Accept()
	if err != nil {
		return err
	}

	if h.config.logs {
		log.Printf("server accepted connection from: '%s'", conn.RemoteAddr().String())
	}

	go h.connectionHandler(conn, handler)
	return nil
}

func (h *HttpServer) Close() error {
	return h.listener.Close()
}

func (h *HttpServer) connectionHandler(conn net.Conn, handler HttpMessageHandler) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	var lines []string

	// headers
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("error while reading: %s", err)
			break
		}

		if h.config.logs {
			trim := strings.TrimSpace(line)
			log.Printf("read message: '%s'", trim)
		}

		// headers section is over
		if line == "\r\n" {
			break
		}

		lines = append(lines, line)
	}

	verb := strings.Split(lines[0], " ")[0]
	headers := extractHeaders(lines)

	lengthStr := headers["Content-Length"]
	contentLength, err := strconv.Atoi(lengthStr)
	if err != nil {
		log.Printf("could not convert string '%s' to int", lengthStr)
		return
	}

	content := extractContent(contentLength, reader)

	if h.config.logs {
		log.Printf("read content: %s", content)
	}

	handler(verb, content, headers)
}

func extractHeaders(lines []string) map[string]string {
	headers := make(map[string]string)

	// skip first line
	for _, line := range lines[1:] {
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			headers[key] = val
		}
	}

	return headers
}

func extractContent(contentLength int, reader *bufio.Reader) string {
	var content []byte

	for range contentLength {
		byte, err := reader.ReadByte()
		if err != nil {
			log.Printf("could not read byte: %s", err)
			return ""
		}

		content = append(content, byte)
	}

	return string(content[:])
}
