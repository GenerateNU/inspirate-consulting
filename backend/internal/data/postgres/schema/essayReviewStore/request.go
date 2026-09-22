package essayReviewRepository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayReviewRepository) RequestEssayReview(ctx context.Context, input models.RequestEssayReviewRequestBody) (*models.EssayReviewTransaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const lockStudentQuery = `
	SELECT review_balance
	FROM public.student
	WHERE id = $1
	FOR UPDATE
	`

	var reviewBalance int
	if err := tx.QueryRow(ctx, lockStudentQuery, input.StudentID).Scan(&reviewBalance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("student", "id", input.StudentID.String())
		}
		return nil, err
	}

	const openReviewQuery = `
	SELECT id
	FROM public.essay_review_transaction
	WHERE essay_id = $1 AND status = 'open'
	LIMIT 1
	`

	var openReviewID uuid.UUID
	switch err := tx.QueryRow(ctx, openReviewQuery, input.EssayID).Scan(&openReviewID); {
	case err == nil:
		return nil, errs.NewHTTPError(http.StatusConflict, fmt.Errorf(
			"essay already has an open review request (%s)", openReviewID,
		))
	case !errors.Is(err, pgx.ErrNoRows):
		return nil, err
	}

	if reviewBalance < input.Amount {
		return nil, errs.BadRequest(fmt.Sprintf(
			"insufficient review balance: student has %d, requested %d", reviewBalance, input.Amount,
		))
	}

	const debitQuery = `
	UPDATE public.student
	SET review_balance = review_balance - $2, updated_at = NOW()
	WHERE id = $1
	`

	if _, err := tx.Exec(ctx, debitQuery, input.StudentID, input.Amount); err != nil {
		return nil, err
	}

	const insertQuery = `
	INSERT INTO public.essay_review_transaction (
		subtotal, student_id, entry_type, essay_id
	) VALUES (
		$1, $2, 'spend', $3
	)
	RETURNING ` + transactionColumns

	rows, err := tx.Query(ctx, insertQuery, -input.Amount, input.StudentID, input.EssayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	createdTransaction, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.EssayReviewTransaction])
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &createdTransaction, nil
}
