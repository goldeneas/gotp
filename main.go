package main

import (
	"log"

	"github.com/goldeneas/gotp/http"
)

func main() {
	router := http.NewHttpRouter()
	router.Add("get", "/", messageHandler)

	server := http.NewServer(
		router,
		http.SetEnableLogs(true),
	)

	server.Listen(25565)
	defer server.Close()

	for {
		server.Accept()
	}
}

func messageHandler(verb string, headers map[string]string, content string) {
	log.Printf("got %s, %s", verb, content)
}
