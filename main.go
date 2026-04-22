package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default config")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := ":" + port

	// Start TCP server
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	fmt.Printf("🚀 IncDB TCP server running on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("Client connected: %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		// Read command (line-based)
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Client disconnected:", err)
			return
		}

		// Trim newline
		msg = trim(msg)

		log.Println("Received:", msg)

		// Simple command handling
		switch msg {
		case "PING":
			conn.Write([]byte("PONG\n"))
		case "INCR":
			conn.Write([]byte("OK\n"))
		case "GET":
			conn.Write([]byte("0\n"))
		case "RESET":
			conn.Write([]byte("OK\n"))
		case "EXIT":
			conn.Write([]byte("BYE\n"))
			return
		default:
			conn.Write([]byte("UNKNOWN COMMAND\n"))
		}
	}
}

// simple trim (avoid strings pkg for now)
func trim(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '\r' {
		s = s[:len(s)-1]
	}
	return s
}
