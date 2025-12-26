package http

import (
	"bufio"
	"strings"
	"testing"
)

func TestExtractFunctions(t *testing.T) {
	rawInput := "POST /api/v1/resource HTTP/1.1\r\n" +
		"Content-Length: 12\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"hello world!"

	reader := bufio.NewReader(strings.NewReader(rawInput))

	lines := extractLines(reader)
	if len(lines) != 3 {
		t.Fatalf("Expected 3 header lines, got %d", len(lines))
	}

	verb := extractVerb(lines)
	if verb != "POST" {
		t.Errorf("Expected verb POST, got %s", verb)
	}

	path := extractPath(lines)
	if path != "/api/v1/resource" {
		t.Errorf("Expected path /api/v1/resource, got %s", path)
	}

	headers := extractHeaders(lines)
	if headers["Content-Length"] != "12" {
		t.Errorf("Expected Content-Length 12, got %s", headers["Content-Length"])
	}
	if headers["Content-Type"] != "text/plain" {
		t.Errorf("Expected Content-Type text/plain, got %s", headers["Content-Type"])
	}

	content, err := readContent(headers, reader)
	if err != nil {
		t.Fatalf("readContent failed: %v", err)
	}
	if content != "hello world!" {
		t.Errorf("Expected body 'hello world!', got '%s'", content)
	}
}

func TestReadContent_MissingHeader(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("some data"))
	emptyHeaders := make(map[string]string)

	_, err := readContent(emptyHeaders, reader)
	if err == nil {
		t.Error("Expected error for missing Content-Length header, got nil")
	}
}
