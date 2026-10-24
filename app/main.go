package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	// Read the port from an environment variable, so Docker and ECS can
	// configure it without changing the code.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &Server{store: newStore()}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, srv.routes()); err != nil {
		log.Fatal(err)
	}
}

func newStore() Store {
	table := os.Getenv("TABLE_NAME")
	if table == "" {
		log.Print("TABLE_NAME not set, using in-memory store")
		return NewMemoryStore()
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("loading AWS config: %v", err)
	}
	log.Printf("using DynamoDB table %s", table)
	return NewDynamoStore(dynamodb.NewFromConfig(cfg), table)
}
