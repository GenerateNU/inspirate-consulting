package essayReviewRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayReviewRepository) SetReviewBalance(ctx context.Context, id uuid.UUID, reviewBalance int) (*models.Student, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const lockStudentQuery = `
	SELECT id, user_id, year, organization, gpa, review_balance, counselor_id
	FROM public.student
	WHERE id = $1
	FOR UPDATE
	`

	student := &models.Student{}
	err = tx.QueryRow(ctx, lockStudentQuery, id).Scan(
		&student.ID,
		&student.UserID,
		&student.Year,
		&student.Organization,
		&student.GPA,
		&student.ReviewBalance,
		&student.CounselorID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("student", "id", id.String())
		}
		return nil, err
	}
	delta := reviewBalance - student.ReviewBalance
	if delta == 0 {
		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
		return student, nil
	}

	const updateBalanceQuery = `
	UPDATE public.student
	SET review_balance = $2, updated_at = NOW()
	WHERE id = $1
	`

	if _, err := tx.Exec(ctx, updateBalanceQuery, id, reviewBalance); err != nil {
		return nil, err
	}

	const insertAdjustmentQuery = `
	INSERT INTO public.essay_review_transaction (
		subtotal, student_id, entry_type
	) VALUES (
		$1, $2, 'adjustment'
	)
	`

	if _, err := tx.Exec(ctx, insertAdjustmentQuery, delta, id); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	student.ReviewBalance = reviewBalance
	return student, nil
}
