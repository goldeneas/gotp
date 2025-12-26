package http

import (
	"bufio"
	"strconv"
	"strings"
)

func readLines(reader *bufio.Reader) ([]string, error) {
	var lines []string

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			return nil, err
		}

		if line == "\r\n" {
			break
		}

		lines = append(lines, line)
	}

	return lines, nil
}

func extractVerb(lines []string) string {
	return strings.Split(lines[0], " ")[0]
}

func extractPath(lines []string) string {
	return strings.Split(lines[0], " ")[1]
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

func readContent(headers map[string]string, reader *bufio.Reader) (string, error) {
	lengthStr := headers["Content-Length"]
	contentLength, err := strconv.Atoi(lengthStr)

	if err != nil {
		return "", err
	}

	var content []byte

	for range contentLength {
		byte, err := reader.ReadByte()
		if err != nil {
			return "", err
		}

		content = append(content, byte)
	}

	return string(content[:]), nil
}
