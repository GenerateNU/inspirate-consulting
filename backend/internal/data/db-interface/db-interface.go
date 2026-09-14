package dbinterface

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// This QueryInterface abstracts what we can run a query on (a pool or transaction),
// *pgxpool.Pool (a connection pool) and pgx.Tx (a transaction) both have the Query, QueryRow, Exec methods,
// Store methods can take in a QueryInterface object and run any of these methods either on 
// a standalone connection, or in a transaction
	// ex)
	// func (s *Store) GetByID(ctx context.Context, db QueryInterface, ud uuid.UUID) (*models.User, error) {
	//   row := db.QueryRow(ctx, query, id)
	// }
// Now when we use the GetByID function, we can either pass in a transaction or a pool as the "db" parameter

type QueryInterface interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}
