package essayRepository

import (
	"context"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"

	"github.com/google/uuid"
)

func (r *EssayRepository) DeleteEssay(ctx context.Context, essayID uuid.UUID, studentID uuid.UUID) error {
	deleteQuery, err := schema.ReadSQLBaseScript("delete_essay.sql", SqlEssayFiles)
	if err != nil {
		return err
	}
	result, err := r.db.Exec(ctx, deleteQuery, essayID, studentID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errs.NotFound("essay", "essayID", essayID)
	}

	return nil

}
