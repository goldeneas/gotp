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

func extractPathAndQueries(lines []string) (string, map[string]string) {
	fullPath := strings.Split(lines[0], " ")[1]

	parts := strings.SplitN(fullPath, "?", 2)
	path := parts[0]

	if len(parts) == 1 {
		return path, nil
	}

	queries := make(map[string]string)
	queryPairs := strings.SplitSeq(parts[1], "&")
	for pair := range queryPairs {
		queryParts := strings.SplitN(pair, "=", 2)
		key := queryParts[0]

		value := ""

		// some queries may not have a value set
		if len(queryParts) == 2 {
			value = queryParts[1]
		}

		queries[key] = value
	}

	return path, queries
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
