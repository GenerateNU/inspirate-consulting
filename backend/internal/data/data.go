package data

import "github.com/jackc/pgx/v5/pgxpool"

//For each schema, their interfaces are to be defined here

type Repository struct {
	db *pgxpool.Pool
	// For each interface, add a field here
}

// Close closes the database connection pool
func (r *Repository) Close() error {
	r.db.Close()
	return nil
}

// GetDB returns the underlying pgxpool.Pool instance
func (r *Repository) GetDB() *pgxpool.Pool {
	return r.db
}

// NewRepository creates a new Repository instance with the given database pool
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
		// For each interface, add an instance of the interface here
	}
}
