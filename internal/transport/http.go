package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cloudemprise/byteSizeGo-goland-api/internal/todo"
)

// TodoItem represents a JSON structure for receiving or sending a single todo item over HTTP.
type TodoItem struct {
	Item string `json:"item"`
}

// Server represents an HTTP server that handles requests related to todos.
type Server struct {
	mux *http.ServeMux // Multiplexer for routing HTTP requests to handlers.
}

// NewServer creates a new Server instance and sets up route handlers for managing todos.
func NewServer(todoService *todo.Service) *Server {
	mux := http.NewServeMux()

	// Handler for fetching all todos via GET request at "/todo".
	mux.HandleFunc("GET /todo", func(w http.ResponseWriter, r *http.Request) {
		// Retrieve all todos using the service layer.
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

	// Handler for searching todos via GET request at "/search?q=<query>".
	mux.HandleFunc("GET /search", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query().Get("q")

		// If no query is provided in URL parameters, return 400 Bad Request.
		if query == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		// Search for todos matching the query using service layer.
		results, err := todoService.Search(query)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Encode search results into JSON format and send response back to client.
		b, err := json.Marshal(results)
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
func (s *Server) Serve() error {
	return http.ListenAndServe(":8080", s.mux) // Bind server to port 8080 using configured mux routes.
}
