package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Item represents a single row in the 'todo_items' table in the database,
// containing both task description and its status (e.g., TO_BE_STARTED).
type Item struct {
	Task   string
	Status string
}

// DB wraps around pgxpool.Pool to manage connections to a PostgreSQL database.
type DB struct {
	pool *pgxpool.Pool // Connection pool for interacting with PostgreSQL database.
}

// New creates a new DB instance by establishing a connection pool to PostgreSQL
// using provided credentials.
// It returns an error if it fails to connect or ping the database successfully.
func New(user, password, host, port, dbname string) (*DB, error) {
	// Construct a connection-string.
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		dbname,
	)

	// Create a database connection pool.
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ping the database to ensure it's reachable before proceeding further.
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{pool: pool}, nil
}

// Close closes all connections in the pool when shutting down or cleaning up resources.
func (db *DB) Close() {
	db.pool.Close()
}

// InsertItem inserts a new task into the 'todo_items' table in PostgreSQL with its status.
// It returns an error if insertion fails due to any reason (e.g., SQL error).
func (db *DB) InsertItem(ctx context.Context, item Item) error {
	query := `INSERT INTO todo_items (task, status) VALUES ($1, $2)`
	_, err := db.pool.Exec(ctx, query, item.Task, item.Status) // Execute SQL insert statement with task and status values.
	return err
}

// GetAllItems retrieves all tasks from 'todo_items' table in PostgreSQL.
// It returns an array of Items or an error if retrieval fails due to any reason (e.g., SQL error).
func (db *DB) GetAllItems(ctx context.Context) ([]Item, error) {
	query := `SELECT task, status FROM todo_items`

	rows, err := db.pool.Query(ctx, query) // Execute SQL select statement to fetch all tasks and statuses from 'todo_items'.
	if err != nil {
		return nil, err
	}

	defer rows.Close() // Ensure rows are closed once processing is complete.

	var items []Item

	// Iterate over each row returned by the query result set and scan values into Item structs.
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.Task, &item.Status); err != nil {
			return nil, err
		}

		items = append(items, item) // Append each scanned item into result slice 'items'.
	}

	// Check for any errors encountered during iteration over rows (e.g., network issues).
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil // Return all fetched items or an empty slice in case of no results found.
}
