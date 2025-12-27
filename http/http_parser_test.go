package http

import (
	"bufio"
	"strings"
	"testing"
)

func TestExtractFunctions(t *testing.T) {
	rawInput := "POST /api/v1/resource?id=123&mode=dark HTTP/1.1\r\n" +
		"Content-Length: 12\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"hello world!"

	reader := bufio.NewReader(strings.NewReader(rawInput))

	lines, err := readLines(reader)
	if err != nil || len(lines) != 3 {
		t.Fatalf("Expected 3 header lines, got %d", len(lines))
	}

	verb := extractVerb(lines)
	if verb != "POST" {
		t.Errorf("Expected verb POST, got %s", verb)
	}

	path, queries := extractPathAndQueries(lines)

	if path != "/api/v1/resource" {
		t.Errorf("Expected path /api/v1/resource, got %s", path)
	}

	if queries["id"] != "123" {
		t.Errorf("Expected query id=123, got %s", queries["id"])
	}
	if queries["mode"] != "dark" {
		t.Errorf("Expected query mode=dark, got %s", queries["mode"])
	}

	headers := extractHeaders(lines)
	if headers["Content-Length"] != "12" {
		t.Errorf("Expected Content-Length 12, got %s", headers["Content-Length"])
	}

	content, err := readContent(headers, reader)
	if err != nil {
		t.Fatalf("readContent failed: %v", err)
	}
	if content != "hello world!" {
		t.Errorf("Expected body 'hello world!', got '%s'", content)
	}
}

func TestParserEdgeCases(t *testing.T) {
	t.Run("Missing Content-Length", func(t *testing.T) {
		raw := "POST / HTTP/1.1\r\nHost: localhost\r\n\r\nSecretBody"
		reader := bufio.NewReader(strings.NewReader(raw))
		lines, _ := readLines(reader)
		headers := extractHeaders(lines)

		_, err := readContent(headers, reader)
		if err == nil {
			t.Error("Expected error when Content-Length is missing for body reading")
		}
	})
}
