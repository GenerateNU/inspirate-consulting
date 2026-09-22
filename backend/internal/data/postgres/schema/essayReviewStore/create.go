package essayReviewRepository

import (
	"context"

	"github.com/google/uuid"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	"inspirate-consulting/internal/models"
)

func (r *EssayReviewRepository) InsertSpend(ctx context.Context, db dbinterface.QueryInterface, studentID, essayID uuid.UUID, amount int) (*models.EssayReviewTransaction, error) {
	const insertQuery = `
	INSERT INTO public.essay_review_transaction (
		subtotal, student_id, entry_type, essay_id
	) VALUES (
		$1, $2, 'spend', $3
	)
	RETURNING ` + transactionColumns

	return collectOne(ctx, db, insertQuery, -amount, studentID, essayID)
}

func (r *EssayReviewRepository) InsertRefund(ctx context.Context, db dbinterface.QueryInterface, charge *models.EssayReviewTransaction) (*models.EssayReviewTransaction, error) {
	const insertQuery = `
	INSERT INTO public.essay_review_transaction (
		subtotal, student_id, entry_type, essay_id
	) VALUES (
		$1, $2, 'refund', $3
	)
	RETURNING ` + transactionColumns

	return collectOne(ctx, db, insertQuery, -charge.Subtotal, charge.StudentID, charge.EssayID)
}

func (r *EssayReviewRepository) InsertAdjustment(ctx context.Context, db dbinterface.QueryInterface, studentID uuid.UUID, delta int) error {
	const insertQuery = `
	INSERT INTO public.essay_review_transaction (
		subtotal, student_id, entry_type
	) VALUES (
		$1, $2, 'adjustment'
	)
	`

	_, err := db.Exec(ctx, insertQuery, delta, studentID)
	return err
}
