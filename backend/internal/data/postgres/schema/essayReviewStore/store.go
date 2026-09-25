package essayReviewRepository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	"inspirate-consulting/internal/models"
)

type EssayReviewRepository struct {
	db *pgxpool.Pool
}

func NewEssayReviewRepository(db *pgxpool.Pool) *EssayReviewRepository {
	return &EssayReviewRepository{db: db}
}

// DB exposes the pool for single-statement reads that need no coordination.
func (r *EssayReviewRepository) DB() dbinterface.QueryInterface {
	return r.db
}

// WithTx runs fn inside one database transaction, committing only if fn
// returns nil. Callers compose queries without importing pgx.
func (r *EssayReviewRepository) WithTx(ctx context.Context, fn func(db dbinterface.QueryInterface) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

const transactionColumns = `id, subtotal, student_id, entry_type, essay_id, completed_at, refund, status, actor_id, created_at, updated_at`

func collectOne(ctx context.Context, db dbinterface.QueryInterface, sql string, args ...any) (*models.EssayReviewTransaction, error) {
	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transaction, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.EssayReviewTransaction])
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
