package essayReviewRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayReviewRepository) GetEssayReviewStatus(ctx context.Context, essayID uuid.UUID) (*models.EssayReviewStatus, error) {

	const selectQuery = `
	SELECT id AS transaction_id, essay_id, student_id, status, created_at AS requested_at, completed_at
	FROM public.essay_review_transaction
	WHERE essay_id = $1 AND entry_type = 'spend'
	ORDER BY created_at DESC, id DESC
	LIMIT 1
	`

	rows, err := r.db.Query(ctx, selectQuery, essayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	status, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.EssayReviewStatus])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("essay review", "essay_id", essayID.String())
		}
		return nil, err
	}

	return &status, nil
}
