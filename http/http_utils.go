package http

import (
	"fmt"
	"net"
)

const HTTP_VERSION = "1.1"

func Send(status string, body string, headers map[string]string, conn net.Conn) {
	fmt.Fprintf(conn, "HTTP/%s %s\r\n", HTTP_VERSION, status)
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))

	for key, value := range headers {
		fmt.Fprintf(conn, "%s: %s\r\n", key, value)
	}

	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, "%s", body)
}

func Serve(file []byte, mime string, conn net.Conn) {
	content := string(file[:])

	Send("200 OK", content, map[string]string{
		"Content-Type": fmt.Sprintf("%s", mime),
	}, conn)
}
