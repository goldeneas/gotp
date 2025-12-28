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

	err := server.Listen(25565)
	if err != nil {
		log.Fatalf("error while trying to listen: %s", err)
	}

	defer server.Close()

	for {
		err := server.Accept()
		if err != nil {
			log.Printf("error while accepting new connection: %s", err)
		}
	}
}

func messageHandler(req *http.HTTPRequest, conn net.Conn) {
	log.Printf("got %s, %s", req.Verb(), req.Content())
	http.Send("200 OK", "Welcome!", nil, conn)
}
