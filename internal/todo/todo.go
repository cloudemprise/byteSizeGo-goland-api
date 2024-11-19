package todo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudemprise/byteSizeGo-goland-api/internal/db"
)

// Item represents a single todo item with a task description and its status.
// This struct is used to model data at the application level (business logic).
type Item struct {
	Task   string // Description of the task (e.g., "Buy groceries").
	Status string // Status of the task (e.g., "TO_BE_STARTED", "COMPLETED").
}

// Manager is an interface that abstracts database operations for inserting and retrieving todo items.
// This abstraction allows easier testing by enabling mock implementations of this interface.
type Manager interface {
	InsertItem(ctx context.Context, item db.Item) error // Inserts a new item into the database.
	GetAllItems(ctx context.Context) ([]db.Item, error) // Retrieves all items from the database.
}

// Service provides methods to interact with todo items in the database.
// It acts as a layer between HTTP handlers and database operations, encapsulating business logic.
type Service struct {
	db Manager // Database connection for storing and retrieving todo items.
}

// NewService creates a new Service instance with the given database connection.
// This function is called in `main.go` to initialize a service with access to DB operations.
func NewService(db Manager) *Service {
	return &Service{
		db: db,
	}
}

// Add inserts a new todo item into the database if it doesn't already exist.
// It checks for duplicates by comparing against existing tasks in the DB before inserting.
func (s *Service) Add(todo string) error {
	// Retrieve all existing todo items from the database to check for duplicates.
	items, err := s.GetAll()

	if err != nil {
		return fmt.Errorf("could not read from db: %w", err) // Return an error if retrieval fails (wrapped with context).
	}

	// Check if the new todo item already exists in the list (case-sensitive comparison).
	for _, t := range items {
		if t.Task == todo {
			return errors.New("todo is not unique") // Return an error if there's already a task with this name.
		}
	}

	// Insert the new todo item into the database with a default status of "TO_BE_STARTED".
	if err = s.db.InsertItem(context.Background(), db.Item{
		Task:   todo,
		Status: "TO_BE_STARTED",
	}); err != nil {
		return fmt.Errorf("could not insert item: %w", err) // Return an error if insertion fails (wrapped with context).
	}

	return nil // Return nil if insertion succeeds without issues.
}

// GetAll retrieves all todo items from the database.
// It converts `db.Item` structs from DB queries into `todo.Item` structs used at this layer.
func (s *Service) GetAll() ([]Item, error) {
	var results []Item

	// Fetch all items from the database using Manager's GetAllItems method.
	items, err := s.db.GetAllItems(context.Background())

	if err != nil {
		return nil, fmt.Errorf("failed to read from db: %w", err) // Return an error if retrieval fails (wrapped with context).
	}

	// Convert each `db.Item` into `todo.Item` before returning them to callers (e.g., HTTP handlers).
	for _, item := range items {
		results = append(results, Item{
			Task:   item.Task,
			Status: item.Status,
		})
	}

	return results, nil // Return all retrieved items or an empty slice in case of no results found.
}

// Search looks for todo items that contain the query string (case-insensitive).
// This method filters tasks by checking whether they contain a substring matching `query`.
func (s *Service) Search(query string) ([]string, error) {

	items, err := s.GetAll() // Retrieve all tasks from DB first since search is done in-memory.

	if err != nil {
		return nil, fmt.Errorf("could not read from db: %w", err)
	}

	var results []string

	// Iterate over each task and check if it contains the query string (case-insensitive).
	for _, todo := range items {
		if strings.Contains(strings.ToLower(todo.Task), strings.ToLower(query)) {
			results = append(results, todo.Task) // Append matching tasks to result slice.
		}
	}

	return results, nil // Return matching tasks or an empty slice if none match.
}
