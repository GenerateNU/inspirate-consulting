package essayReviewRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// LockStudent reads a student and holds a row lock until the surrounding
// transaction ends.
func (r *EssayReviewRepository) LockStudent(ctx context.Context, db dbinterface.QueryInterface, id uuid.UUID) (*models.Student, error) {
	lockQuery, err := schema.ReadSQLBaseScript("get_essay_review.sql", SqlEssayReviewFiles)
	if err != nil {
		return nil, err
	}

	student := &models.Student{}
	err = db.QueryRow(ctx, lockQuery, id).Scan(
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
	lockQuery, err := schema.ReadSQLBaseScript("lock_essay_review_transaction.sql", SqlEssayReviewFiles)
	if err != nil {
		return nil, err
	}

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
	selectQuery, err := schema.ReadSQLBaseScript("find_open_essay_review.sql", SqlEssayReviewFiles)
	if err != nil {
		return nil, err
	}

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
	selectQuery, err := schema.ReadSQLBaseScript("get_essay_review_status.sql", SqlEssayReviewFiles)
	if err != nil {
		return nil, err
	}

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
