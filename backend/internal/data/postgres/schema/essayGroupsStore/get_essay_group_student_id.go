package essayGroupRepository

import (
	"context"
	"errors"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *EssayGroupRepository) GetEssayGroupStudentID(ctx context.Context, groupID uuid.UUID) (uuid.UUID, error) {
	query, err := schema.ReadSQLBaseScript("get_essay_group_student_id.sql", SqlEssayGroupFiles)
	if err != nil {
		e := errs.InternalServerError("Failed to read base query: ", err.Error())
		return uuid.Nil, &e
	}

	var studentID uuid.UUID
	if err := r.db.QueryRow(ctx, query, groupID).Scan(&studentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errs.NotFound("essay group", "id", groupID.String())
		}
		return uuid.Nil, err
	}

	return studentID, nil

}
