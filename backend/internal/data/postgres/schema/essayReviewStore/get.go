package essayReviewRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayReviewRepository) GetEssayReview(ctx context.Context, id uuid.UUID) (*models.EssayReviewTransaction, error) {

	const selectQuery = `
	SELECT ` + transactionColumns + `
	FROM public.essay_review_transaction
	WHERE id = $1
	`

	rows, err := r.db.Query(ctx, selectQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transaction, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.EssayReviewTransaction])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("essay review", "id", id.String())
		}
		return nil, err
	}

	return &transaction, nil
}
