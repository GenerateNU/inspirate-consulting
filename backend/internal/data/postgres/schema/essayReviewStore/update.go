package essayReviewRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayReviewRepository) SetStudentBalance(ctx context.Context, db dbinterface.QueryInterface, id uuid.UUID, reviewBalance int) error {
	const updateQuery = `
	UPDATE public.student
	SET review_balance = $2, updated_at = NOW()
	WHERE id = $1
	`

	tag, err := db.Exec(ctx, updateQuery, id, reviewBalance)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errs.NotFound("student", "id", id.String())
	}

	return nil
}

func (r *EssayReviewRepository) AdjustStudentBalance(ctx context.Context, db dbinterface.QueryInterface, id uuid.UUID, delta int) error {
	const adjustQuery = `
	UPDATE public.student
	SET review_balance = review_balance + $2, updated_at = NOW()
	WHERE id = $1
	`

	tag, err := db.Exec(ctx, adjustQuery, id, delta)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errs.NotFound("student", "id", id.String())
	}

	return nil
}

// LinkRefund points a charge at the row that reversed it.
func (r *EssayReviewRepository) LinkRefund(ctx context.Context, db dbinterface.QueryInterface, chargeID, reversalID uuid.UUID) (*models.EssayReviewTransaction, error) {
	const updateQuery = `
	UPDATE public.essay_review_transaction
	SET refund = $2, updated_at = NOW()
	WHERE id = $1
	RETURNING ` + transactionColumns

	transaction, err := collectOne(ctx, db, updateQuery, chargeID, reversalID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("essay review", "id", chargeID.String())
		}
		return nil, err
	}

	return transaction, nil
}

func (r *EssayReviewRepository) MarkCompleted(ctx context.Context, db dbinterface.QueryInterface, id uuid.UUID) (*models.EssayReviewTransaction, error) {
	const updateQuery = `
	UPDATE public.essay_review_transaction
	SET completed_at = NOW(), updated_at = NOW()
	WHERE id = $1
	RETURNING ` + transactionColumns

	transaction, err := collectOne(ctx, db, updateQuery, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("essay review", "id", id.String())
		}
		return nil, err
	}

	return transaction, nil
}
