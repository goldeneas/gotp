package main

import (
	"log"

	"github.com/goldeneas/gotp/http"
)

func main() {
	server := http.NewServer(
		http.SetEnableLogs(true),
	)

	server.Listen(25565)

	for {
		server.Accept(messageHandler)
	}
}

func messageHandler(verb string, content string, headers map[string]string) {
	log.Printf("got %s, %s", verb, content)
}
