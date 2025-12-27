package http

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type HttpFileServer struct {
	cache map[string][]byte
	mtx   sync.RWMutex
}

func NewHttpFileServer() *HttpFileServer {
	return &HttpFileServer{
		cache: make(map[string][]byte),
	}
}

func (h *HttpFileServer) Get(path string) ([]byte, error) {
	h.mtx.RLock()
	file, exists := h.loadFromCache(path)
	h.mtx.RUnlock()

	if exists {
		return file, nil
	}

	file, err := h.loadFromDisk(path)
	if err != nil {
		return nil, err
	}

	h.mtx.Lock()
	h.cachePut(path, file)
	h.mtx.Unlock()

	return file, nil
}

func (h *HttpFileServer) Invalidate(paths ...string) {
	h.mtx.Lock()

	for _, path := range paths {
		delete(h.cache, path)
		log.Printf("invalidating %s", path)
	}

	h.mtx.Unlock()
}

func (h *HttpFileServer) WatchDir(path string) {
	state := make(map[string]int64)

	for {
		paths, err := h.scanChanges(path, state)
		if err != nil {
			fmt.Printf("error while watching files: %s", err)
		}

		if len(paths) > 0 {
			h.Invalidate(paths...)
		}

		time.Sleep(5 * time.Second)
	}
}

func (h *HttpFileServer) cachePut(key string, file []byte) {
	h.cache[key] = file
}

func (h *HttpFileServer) loadFromCache(key string) ([]byte, bool) {
	file, ok := h.cache[key]
	return file, ok
}

func (h *HttpFileServer) scanChanges(path string, state map[string]int64) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var changed []string
	seen := make(map[string]bool)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := info.Name()
		curr := info.ModTime().UnixMicro()

		fullPath := filepath.Join(path, name)

		prev, exists := state[fullPath]
		if exists && prev != curr {
			changed = append(changed, fullPath)
		}

		state[fullPath] = curr
		seen[fullPath] = true
	}

	for fullPath := range state {
		if seen[fullPath] {
			continue
		}

		changed = append(changed, fullPath)
		delete(state, fullPath)
	}

	return changed, nil
}

func (h *HttpFileServer) loadFromDisk(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}
