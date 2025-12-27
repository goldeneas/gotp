package http

import (
	"os"
	"path/filepath"
)

type HttpFileServer struct {
	cache map[string][]byte
}

func NewHttpFileServer() *HttpFileServer {
	return &HttpFileServer{
		cache: make(map[string][]byte),
	}
}

func (h *HttpFileServer) Get(prefix string, path string) ([]byte, error) {
	if h.cacheContains(path) {
		return h.loadFromCache(path), nil
	}

	fullPath := filepath.Join(prefix, path)
	file, err := h.loadFromDisk(fullPath)
	if err != nil {
		return nil, err
	}

	h.cachePut(path, file)
	return file, nil
}

func (h *HttpFileServer) cachePut(path string, file []byte) {
	h.cache[path] = file
}

func (h *HttpFileServer) cacheContains(path string) bool {
	_, contains := h.cache[path]
	return contains
}

func (h *HttpFileServer) loadFromCache(path string) []byte {
	return h.cache[path]
}

func (h *HttpFileServer) loadFromDisk(fullPath string) ([]byte, error) {
	file, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
