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

func (r *EssayReviewRepository) RefundEssayReview(ctx context.Context, id uuid.UUID) (*models.RefundEssayReviewResponseBody, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const lockTransactionQuery = `
	SELECT ` + transactionColumns + `
	FROM public.essay_review_transaction
	WHERE id = $1
	FOR UPDATE
	`

	lockRows, err := tx.Query(ctx, lockTransactionQuery, id)
	if err != nil {
		return nil, err
	}
	defer lockRows.Close()

	originalTransaction, err := pgx.CollectExactlyOneRow(lockRows, pgx.RowToStructByName[models.EssayReviewTransaction])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("essay review", "id", id.String())
		}
		return nil, err
	}

	switch {
	case originalTransaction.Status == nil:
		return nil, errs.BadRequest(fmt.Sprintf(
			"only a spend can be refunded, this entry is a %s", originalTransaction.EntryType,
		))
	case *originalTransaction.Status == models.ReviewStatusCompleted:
		return nil, errs.NewHTTPError(http.StatusConflict, errors.New(
			"review has already been completed and cannot be refunded",
		))
	case *originalTransaction.Status == models.ReviewStatusRefunded:
		return nil, errs.NewHTTPError(http.StatusConflict, errors.New(
			"review has already been refunded",
		))
	}

	const lockStudentQuery = `
	SELECT id
	FROM public.student
	WHERE id = $1
	FOR UPDATE
	`

	var lockedStudentID uuid.UUID
	if err := tx.QueryRow(ctx, lockStudentQuery, originalTransaction.StudentID).Scan(&lockedStudentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("student", "id", originalTransaction.StudentID.String())
		}
		return nil, err
	}

	const insertRefundQuery = `
	INSERT INTO public.essay_review_transaction (
		subtotal, student_id, entry_type, essay_id
	) VALUES (
		$1, $2, 'refund', $3
	)
	RETURNING ` + transactionColumns

	insertRows, err := tx.Query(
		ctx,
		insertRefundQuery,
		-originalTransaction.Subtotal,
		originalTransaction.StudentID,
		originalTransaction.EssayID,
	)
	if err != nil {
		return nil, err
	}
	defer insertRows.Close()

	refundTransaction, err := pgx.CollectExactlyOneRow(insertRows, pgx.RowToStructByName[models.EssayReviewTransaction])
	if err != nil {
		return nil, err
	}

	const linkRefundQuery = `
	UPDATE public.essay_review_transaction
	SET refund = $2, updated_at = NOW()
	WHERE id = $1 AND status = 'open'
	RETURNING ` + transactionColumns

	linkRows, err := tx.Query(ctx, linkRefundQuery, id, refundTransaction.ID)
	if err != nil {
		return nil, err
	}
	defer linkRows.Close()

	refundedTransaction, err := pgx.CollectExactlyOneRow(linkRows, pgx.RowToStructByName[models.EssayReviewTransaction])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewHTTPError(http.StatusConflict, errors.New(
				"review is no longer refundable",
			))
		}
		return nil, err
	}

	const creditQuery = `
	UPDATE public.student
	SET review_balance = review_balance + $2, updated_at = NOW()
	WHERE id = $1
	`

	creditTag, err := tx.Exec(ctx, creditQuery, originalTransaction.StudentID, refundTransaction.Subtotal)
	if err != nil {
		return nil, err
	}
	if creditTag.RowsAffected() != 1 {
		return nil, errs.InternalServerError(fmt.Sprintf(
			"expected to credit exactly one student, credited %d", creditTag.RowsAffected(),
		))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.RefundEssayReviewResponseBody{
		Transaction: refundedTransaction,
		Refund:      refundTransaction,
	}, nil
}
