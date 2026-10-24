package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Read the port from an environment variable, so Docker and ECS can
	// configure it without changing the code.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &Server{store: NewMemoryStore()}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, srv.routes()); err != nil {
		log.Fatal(err)
	}
}
