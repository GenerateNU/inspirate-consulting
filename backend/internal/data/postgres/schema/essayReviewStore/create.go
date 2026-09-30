package essayReviewRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayReviewRepository) InsertSpend(ctx context.Context, db dbinterface.QueryInterface, studentID, essayID uuid.UUID, amount int) (*models.EssayReviewTransaction, error) {
	insertQuery, err := schema.ReadSQLBaseScript("create_essay_review.sql", SqlEssayReviewFiles)

	if err != nil {
		fmt.Println("[db] ReadSQLBaseScript error: ", err)
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, err
	}

	return collectOne(ctx, db, insertQuery, -amount, studentID, essayID)
}

func (r *EssayReviewRepository) InsertRefund(ctx context.Context, db dbinterface.QueryInterface, charge *models.EssayReviewTransaction) (*models.EssayReviewTransaction, error) {
	insertQuery, err := schema.ReadSQLBaseScript("insert_essay_review_refund.sql", SqlEssayReviewFiles)
	if err != nil {
		return nil, err
	}

	return collectOne(ctx, db, insertQuery, -charge.Subtotal, charge.StudentID, charge.EssayID)
}

func (r *EssayReviewRepository) InsertAdjustment(ctx context.Context, db dbinterface.QueryInterface, studentID uuid.UUID, delta int) error {
	insertQuery, err := schema.ReadSQLBaseScript("insert_essay_review_adjustment.sql", SqlEssayReviewFiles)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, insertQuery, delta, studentID)
	return err
}
