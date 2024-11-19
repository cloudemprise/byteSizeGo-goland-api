package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cloudemprise/byteSizeGo-goland-api/internal/todo"
)

// TodoItem represents a JSON structure for receiving or sending a single todo item over HTTP.
// It contains only one field 'Item', which is expected in JSON format when adding todos via POST request.
type TodoItem struct {
	Item string `json:"item"` // The task description sent/received as JSON data via HTTP requests/responses
}

// Server represents an HTTP server that handles requests related to todos.
// It uses http.ServeMux as its multiplexer for routing incoming HTTP requests to appropriate handlers.
type Server struct {
	mux *http.ServeMux // Multiplexer for routing HTTP requests to handlers based on URL paths and methods
}

// NewServer creates a new Server instance and sets up route handlers for managing todos.
// It takes in a reference to `todo.Service`, which contains business logic for adding/retrieving/searching todos.
func NewServer(todoService *todo.Service) *Server {
	mux := http.NewServeMux()

	// Handler for fetching all todos via GET request at "/todo".
	mux.HandleFunc("GET /todo", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve all todos using service layer's GetAll method.
		items, err := todoService.GetAll()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Encode todos into JSON format and send response back to client.
		b, err := json.Marshal(items)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			log.Fatal(err)
		}
	})

	// Handler for adding a new todo via POST request at "/todo".
	mux.HandleFunc("POST /todo", func(writer http.ResponseWriter, request *http.Request) {
		var t TodoItem

		// Decode incoming JSON request into TodoItem struct.
		err := json.NewDecoder(request.Body).Decode(&t)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		// Add new todo item using service layer. If there's an error (e.g., duplicate), return 400 Bad Request.
		err = todoService.Add(t.Item)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.WriteHeader(http.StatusCreated) // Return 201 Created on success.
	})

	// Handler for searching todos via GET request at "/search?q=".
	mux.HandleFunc("GET /search", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query().Get("q")

		// If no query is provided in URL parameters, return 400 Bad Request.
		if query == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		// Search for todos matching the query using service layer's Search method.
		results, err := todoService.Search(query)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(results) // Encode search results into JSON format and send response back to client.
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = writer.Write(b)
		if err != nil {
			log.Println(err)
			return
		}
	})

	return &Server{mux: mux} // Return server instance with configured routes.

}

// Serve starts the HTTP server on port 8080 and listens for incoming requests.
// It binds the server to port ":8080" using ServeMux for handling routes defined earlier.
func (s *Server) Serve() error {
	return http.ListenAndServe(":8080", s.mux)
}
