package essayReviewRepository

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayReviewRepository) CompleteEssayReview(ctx context.Context, id uuid.UUID) (*models.EssayReviewTransaction, error) {

	const completeQuery = `
	UPDATE public.essay_review_transaction
	SET completed_at = NOW(), updated_at = NOW()
	WHERE id = $1 AND status = 'open'
	RETURNING ` + transactionColumns

	rows, err := r.db.Query(ctx, completeQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	completedTransaction, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.EssayReviewTransaction])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Re-read to tell a missing review (404) from one that is no
			// longer open (409).
			existingTransaction, getErr := r.GetEssayReview(ctx, id)
			if getErr != nil {
				return nil, getErr
			}
			switch {
			case existingTransaction.Status == nil:
				return nil, errs.BadRequest("only a spend can be completed")
			case *existingTransaction.Status == models.ReviewStatusRefunded:
				return nil, errs.NewHTTPError(http.StatusConflict, errors.New(
					"review has already been refunded and cannot be completed",
				))
			default:
				return nil, errs.NewHTTPError(http.StatusConflict, errors.New(
					"review has already been completed",
				))
			}
		}
		return nil, err
	}

	return &completedTransaction, nil
}
