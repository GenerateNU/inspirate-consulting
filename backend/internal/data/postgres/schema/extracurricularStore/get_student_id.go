package extracurricularRepository

import (
	"context"
	"errors"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *ExtracurricularRepository) GetExtracurricularStudentID(ctx context.Context, id int64) (uuid.UUID, error) {
	query, err := schema.ReadSQLBaseScript("get_student_id.sql", SqlExtracurricularFiles)
	if err != nil {
		e := errs.InternalServerError("Failed to read base query: ", err.Error())
		return uuid.Nil, &e
	}

	var studentID uuid.UUID
	if err := r.db.QueryRow(ctx, query, id).Scan(&studentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errs.NotFound("extracurricular", "id", strconv.FormatInt(id, 10))
		}
		return uuid.Nil, err
	}
	return studentID, nil

}
