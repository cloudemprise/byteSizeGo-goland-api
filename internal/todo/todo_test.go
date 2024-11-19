package todo_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudemprise/byteSizeGo-goland-api/internal/db"
	"github.com/cloudemprise/byteSizeGo-goland-api/internal/todo"
)

type MockDB struct {
	items []db.Item
}

func (m *MockDB) InsertItem(_ context.Context, item db.Item) error {
	m.items = append(m.items, item)
	return nil
}

func (m *MockDB) GetAllItems(_ context.Context) ([]db.Item, error) {
	return m.items, nil
}

func TestService_Search(t *testing.T) {

	tests := []struct {
		name            string
		toDosToAdd      []string
		query           string
		expectedResults []string
	}{
		{
			name:            "given a todo of shop and a search of sh, i should get shop back",
			toDosToAdd:      []string{"shop"},
			query:           "sh",
			expectedResults: []string{"shop"},
		},
		{
			name:            "still returns shop, even if the case doesn't match",
			toDosToAdd:      []string{"Shopping"},
			query:           "sh",
			expectedResults: []string{"Shopping"},
		},
		{
			name:            "spaces",
			toDosToAdd:      []string{"go Shopping"},
			query:           "go",
			expectedResults: []string{"go Shopping"},
		},
		{
			name:            "space at start of word",
			toDosToAdd:      []string{" Space at beginning"},
			query:           "space",
			expectedResults: []string{" Space at beginning"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDB{}          // Instantiate a mock database.
			s := todo.NewService(m) // Create a new todo service.

			// Populate the todo store with elements.
			for _, toAdd := range tt.toDosToAdd {
				err := s.Add(toAdd)
				if err != nil {
					t.Error(err)
				}
			}

			// Perform a lookup within the store.
			got, err := s.Search(tt.query)
			if err != nil {
				t.Error(err)
			}

			// Check that we get back what we expect.
			if !reflect.DeepEqual(got, tt.expectedResults) {
				t.Errorf("Search() = %v, want %v", got, tt.expectedResults)
			}
		})
	}
}
