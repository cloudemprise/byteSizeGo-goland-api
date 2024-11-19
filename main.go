package main

import (
	"log"

	"github.com/cloudemprise/byteSizeGo-goland-api/internal/db"
	"github.com/cloudemprise/byteSizeGo-goland-api/internal/todo"
	"github.com/cloudemprise/byteSizeGo-goland-api/internal/transport"
)

func main() {
	// Initialize a new database connection with the provided credentials and settings.
	// These credentials are likely specific to the local environment setup.
	dbStore, err := db.New(`bytesizego`, "pa55word", "localhost", "5432", "todos")

	// If there's an error initializing the database connection, log it and exit the program.
	if err != nil {
		log.Fatal(err) // log.Fatal logs the error and then terminates the program.
	}

	// Create a new todo service that will interact with the database.
	// The service layer abstracts business logic related to todos.
	service := todo.NewService(dbStore)

	// Initialize a new HTTP server that will handle API requests for the todo service.
	// The transport layer handles routing and HTTP request/response handling.
	server := transport.NewServer(service)

	// Start the server and listen for incoming requests on a predefined port (likely 8080).
	// If there's an error starting or running the server, log it and exit.
	if err := server.Serve(); err != nil {
		log.Fatal(err) // log.Fatal ensures that any critical error in serving is logged before termination.
	}
}
