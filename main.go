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

func messageHandler(req *http.HTTPRequest, conn net.Conn) {
	log.Printf("got %s, %s", req.Verb(), req.Content())
	http.Send("200", "Welcome!", nil, conn)
}
