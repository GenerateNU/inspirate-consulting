package counselorRepository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *CounselorRepository) ListPendingReviews(ctx context.Context, counselorID string) ([]models.PendingReview, error) {
	listQuery, err := schema.ReadSQLBaseScript("list_pending_reviews_by_counselor.sql", SqlCounselorFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, listQuery, counselorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.PendingReview])
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
