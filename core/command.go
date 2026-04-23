package core

import (
	"bufio"
	"log"
	"net"
)

func HandleCommand(conn net.Conn) {
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
