package essayReviewRepository

import "github.com/jackc/pgx/v5/pgxpool"

type EssayReviewRepository struct {
	db *pgxpool.Pool
}

func NewEssayReviewRepository(db *pgxpool.Pool) *EssayReviewRepository {
	return &EssayReviewRepository{db: db}
}

const transactionColumns = `id, subtotal, student_id, entry_type, essay_id, completed_at, refund, status, actor_id, created_at, updated_at`
