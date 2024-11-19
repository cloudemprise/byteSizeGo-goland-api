package todo_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudemprise/byteSizeGo-goland-api/internal/db"
	"github.com/cloudemprise/byteSizeGo-goland-api/internal/todo"
)

// MockDB is a mock implementation of the db.Manager interface.
// It simulates a database by storing items in memory (in the `items` slice).
type MockDB struct {
	items []db.Item // In-memory storage for todo items
}

// InsertItem simulates inserting an item into the database.
// It appends the given item to the `items` slice.
func (m *MockDB) InsertItem(_ context.Context, item db.Item) error {
	m.items = append(m.items, item)
	return nil // Always returns nil since this is a mock and no real error handling is needed
}

// GetAllItems simulates retrieving all items from the database.
// It returns all items stored in memory.
func (m *MockDB) GetAllItems(_ context.Context) ([]db.Item, error) {
	return m.items, nil // Returns all items without any error
}

// TestService_Search tests the Search function of the todo.Service.
// It defines multiple test cases to ensure that searching for todo items works correctly.
func TestService_Search(t *testing.T) {
	// Define test cases as a slice of anonymous structs.
	tests := []struct {
		name            string   // Name of the test case
		toDosToAdd      []string // List of todos to add before performing the search
		query           string   // The query string used for searching
		expectedResults []string // Expected results after performing the search
	}{
		{
			name:            "given a todo of shop and a search of sh, I should get shop back",
			toDosToAdd:      []string{"shop"}, // Adding "shop" as a todo item
			query:           "sh",             // Searching for "sh"
			expectedResults: []string{"shop"}, // Expecting "shop" as the result
		},
		{
			name:            "still returns shop, even if the case doesn't match",
			toDosToAdd:      []string{"Shopping"}, // Adding "Shopping" (with capital S)
			query:           "sh",                 // Searching for "sh" (lowercase)
			expectedResults: []string{"Shopping"}, // Expecting "Shopping" as result (case-insensitive search)
		},
		{
			name:            "spaces",
			toDosToAdd:      []string{"go Shopping"}, // Adding "go Shopping"
			query:           "go",                    // Searching for "go"
			expectedResults: []string{"go Shopping"}, // Expecting "go Shopping" as result
		},
		{
			name:            "space at start of word",
			toDosToAdd:      []string{" Space at beginning"}, // Adding todo with leading space
			query:           "space",                         // Searching for "space"
			expectedResults: []string{" Space at beginning"}, // Expecting result with leading space intact
		},
	}

	// Iterate over each test case.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDB{}          // Create a new MockDB instance (simulating an empty database)
			s := todo.NewService(m) // Create a new todo service using the mock database

			// Add each todo item from the test case to the service.
			for _, toAdd := range tt.toDosToAdd {
				err := s.Add(toAdd) // Add each todo item using the Add method from Service
				if err != nil {
					t.Error(err) // If there's an error while adding, fail the test
				}
			}

			// Perform search using the query from the test case.
			got, err := s.Search(tt.query)
			if err != nil {
				t.Error(err) // If there's an error during search, fail the test
			}

			// Compare actual results with expected results using reflect.DeepEqual.
			if !reflect.DeepEqual(got, tt.expectedResults) {
				t.Errorf("Search() = %v, want %v", got, tt.expectedResults) // Fail if results don't match expected output
			}
		})
	}
}
