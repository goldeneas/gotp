package main

import (
	"log"
	"net"

	"github.com/goldeneas/gotp/http"
)

func main() {
	router := http.NewHttpRouter()
	router.Add("get", "/", messageHandler)
	router.Add("post", "/", messageHandler)

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

func messageHandler(verb string, headers map[string]string, content string, conn net.Conn) {
	log.Printf("got %s, %s", verb, content)
}
