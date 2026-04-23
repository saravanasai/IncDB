package server

import (
	"fmt"
	"log"
	"net"

	"github.com/saravanasai/IncDB/core"
)

func Start(port string) {

	address := ":" + port

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	fmt.Printf("IncDB TCP server running on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go core.HandleCommand(conn)
	}
}
