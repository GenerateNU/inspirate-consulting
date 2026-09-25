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

// LockStudent reads a student and holds a row lock until the surrounding
// transaction ends.
func (r *EssayReviewRepository) LockStudent(ctx context.Context, db dbinterface.QueryInterface, id uuid.UUID) (*models.Student, error) {
	const lockQuery = `
	SELECT id, user_id, year, organization, gpa, review_balance, counselor_id
	FROM public.student
	WHERE id = $1
	FOR UPDATE
	`

	student := &models.Student{}
	err := db.QueryRow(ctx, lockQuery, id).Scan(
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

	return student, nil
}

// LockTransaction reads a transaction row and holds a row lock until the
// surrounding transaction ends.
func (r *EssayReviewRepository) LockTransaction(ctx context.Context, db dbinterface.QueryInterface, id uuid.UUID) (*models.EssayReviewTransaction, error) {
	const lockQuery = `
	SELECT ` + transactionColumns + `
	FROM public.essay_review_transaction
	WHERE id = $1
	FOR UPDATE
	`

	transaction, err := collectOne(ctx, db, lockQuery, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("essay review", "id", id.String())
		}
		return nil, err
	}

	return transaction, nil
}

// FindOpenReviewForEssay returns nil when the essay has no review in flight.
func (r *EssayReviewRepository) FindOpenReviewForEssay(ctx context.Context, db dbinterface.QueryInterface, essayID uuid.UUID) (*uuid.UUID, error) {
	const selectQuery = `
	SELECT id
	FROM public.essay_review_transaction
	WHERE essay_id = $1 AND status = 'open'
	LIMIT 1
	`

	var openReviewID uuid.UUID
	switch err := db.QueryRow(ctx, selectQuery, essayID).Scan(&openReviewID); {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, err
	}

	return &openReviewID, nil
}

// FindEssayReviewStatus returns the most recent spend for an essay.
func (r *EssayReviewRepository) FindEssayReviewStatus(ctx context.Context, db dbinterface.QueryInterface, essayID uuid.UUID) (*models.EssayReviewStatus, error) {
	const selectQuery = `
	SELECT id AS transaction_id, essay_id, student_id, status, created_at AS requested_at, completed_at
	FROM public.essay_review_transaction
	WHERE essay_id = $1 AND entry_type = 'spend'
	ORDER BY created_at DESC, id DESC
	LIMIT 1
	`

	rows, err := db.Query(ctx, selectQuery, essayID)
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
