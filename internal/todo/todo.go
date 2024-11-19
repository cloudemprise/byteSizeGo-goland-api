package todo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudemprise/byteSizeGo-goland-api/internal/db"
)

// Item represents a single todo item with a task description and its status.
type Item struct {
	Task   string
	Status string
}

type Manager interface {
	InsertItem(ctx context.Context, item db.Item) error
	GetAllItems(ctx context.Context) ([]db.Item, error)
}

// Service provides methods to interact with todo items in the database.
type Service struct {
	db Manager // Database connection for storing and retrieving todo items.
}

// NewService creates a new Service instance with the given database connection.
func NewService(db Manager) *Service {
	return &Service{
		db: db,
	}
}

// Add inserts a new todo item into the database if it doesn't already exist.
func (s *Service) Add(todo string) error {
	// Retrieve all existing todo items from the database.
	items, err := s.GetAll()
	if err != nil {
		return fmt.Errorf("could not read from db: %w", err)
	}

	// Check if the new todo item already exists in the list.
	for _, t := range items {
		if t.Task == todo {
			return errors.New("todo is not unique")
		}
	}

	// Insert the new todo item into the database with a default status of "TO_BE_STARTED".
	if err = s.db.InsertItem(context.Background(), db.Item{
		Task:   todo,
		Status: "TO_BE_STARTED",
	}); err != nil {
		return fmt.Errorf("could not insert item: %w", err)
	}

	return nil
}

// GetAll retrieves all todo items from the database.
func (s *Service) GetAll() ([]Item, error) {
	var results []Item

	// Fetch all items from the database.
	items, err := s.db.GetAllItems(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to read from db: %w", err)
	}

	// Convert db.Item to todo.Item and append them to results.
	for _, item := range items {
		results = append(results, Item{
			Task:   item.Task,
			Status: item.Status,
		})
	}

	return results, nil
}

// Search looks for todo items that contain the query string (case-insensitive).
func (s *Service) Search(query string) ([]string, error) {
	items, err := s.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not read from db: %w", err)
	}

	var results []string
	for _, todo := range items {
		if strings.Contains(strings.ToLower(todo.Task), strings.ToLower(query)) {
			results = append(results, todo.Task)
		}
	}
	return results, nil
}
