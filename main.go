package main

import (
	"log"

	"github.com/cloudemprise/byteSizeGo-goland-api/internal/db"
	"github.com/cloudemprise/byteSizeGo-goland-api/internal/todo"
	"github.com/cloudemprise/byteSizeGo-goland-api/internal/transport"
)

func main() {

	// Initialize a new database connection with the provided credentials and settings.
	dbStore, err := db.New(`bytesizego`, "pa55word", "localhost", "5432", "todos")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new todo service that will interact with the database.
	service := todo.NewService(dbStore)

	// Initialize a new HTTP server that will handle API requests for the todo service.
	server := transport.NewServer(service)

	// Start the server and listen for incoming requests. If there's an error, log it and exit.
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
