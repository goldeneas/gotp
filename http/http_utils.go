package http

import (
	"fmt"
	"net"
)

func Send(status string, body string, headers map[string]string, conn net.Conn) {
	fmt.Fprintf(conn, "HTTP/1.1 %s\r\n", status)
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))

	for key, value := range headers {
		fmt.Fprintf(conn, key, value+"\r\n")
	}

	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, "%s", body)
}
