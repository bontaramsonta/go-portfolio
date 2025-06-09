package main

import (
	"log"
	"os"
	"portfolio/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.New()
	
	log.Printf("Starting portfolio server on port %s", port)
	log.Printf("Visit http://localhost:%s to view the site", port)
	
	if err := srv.Start(port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}