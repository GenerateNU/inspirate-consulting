package essayRepository

import (
	"context"
	"errors"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *EssayRepository) GetEssayStudentID(ctx context.Context, essayID uuid.UUID) (uuid.UUID, error) {
	query, err := schema.ReadSQLBaseScript("get_essay_student_id.sql", SqlEssayFiles)
	if err != nil {
		e := errs.InternalServerError("Failed to read base query: ", err.Error())
		return uuid.Nil, &e
	}

	var studentID uuid.UUID
	if err := r.db.QueryRow(ctx, query, essayID).Scan(&studentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errs.NotFound("essay", "id", essayID.String())
		}
		return uuid.Nil, err
	}

	return studentID, nil
}
