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

// updates a student's profile fields, leaving any nil field unchanged
func (r *StudentRepository) UpdateStudent(ctx context.Context, id uuid.UUID, input models.UpdateStudentBody) (*models.Student, error) {
	Student := &models.Student{}

	query, err := schema.ReadSQLBaseScript("update_student.sql", SqlStudentFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	err = r.db.QueryRow(
		ctx,
		query,
		id,
		input.Year,
		input.Organization,
		input.GPA,
		input.Interest,
		input.HighSchool,
		input.StandardizedTestScore,
		input.IntendedMajor,
		input.IntendedMinor,
		input.GeographicPreferences,
	).Scan(
		&Student.ID, &Student.UserID, &Student.Year, &Student.Organization, &Student.GPA,
		&Student.ReviewBalance, &Student.CounselorID, &Student.Interest,
		&Student.HighSchool, &Student.StandardizedTestScore, &Student.IntendedMajor,
		&Student.IntendedMinor, &Student.GeographicPreferences,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("student", "id", id.String())
		}
		return nil, err
	}
	return Student, nil
}
