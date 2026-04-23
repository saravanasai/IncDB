package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/saravanasai/IncDB/server"
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

	server.Start(port)
}
