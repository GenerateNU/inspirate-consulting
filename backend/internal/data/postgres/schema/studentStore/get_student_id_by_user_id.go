package studentRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
)

// GetStudentIDByUserID returns the student ID for a student's user ID.
func (r *StudentRepository) GetStudentIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	query, err := schema.ReadSQLBaseScript("get_student_id_by_user_id.sql", SqlStudentFiles)
	if err != nil {
		return uuid.Nil, err
	}

	var studentID uuid.UUID
	if err := r.db.QueryRow(ctx, query, userID).Scan(&studentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errs.NotFound("student", "user_id", userID)
		}
		return uuid.Nil, err
	}

	return studentID, nil
}
