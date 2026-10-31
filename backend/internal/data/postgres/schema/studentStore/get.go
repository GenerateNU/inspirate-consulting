package studentRepository

import (
	"context"
	"errors"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *StudentRepository) GetStudent(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	Student := &models.Student{}

	query, err := schema.ReadSQLBaseScript("get_student.sql", SqlStudentFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	err = r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(&Student.ID, &Student.UserID, &Student.Year, &Student.Organization, &Student.GPA, &Student.ReviewBalance, &Student.CounselorID, &Student.Interest)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("student", "id", id.String())
		}
		return nil, err
	}
	return Student, nil
}
