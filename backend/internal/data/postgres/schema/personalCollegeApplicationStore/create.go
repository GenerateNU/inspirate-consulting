package personalCollegeApplicationRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// Postgres SQLSTATE codes relevant to this insert.
const (
	pgForeignKeyViolationCode = "23503"
	pgCheckViolationCode      = "23514"
)

func (r *PersonalCollegeApplicationRepository) CreatePersonalCollegeApplication(
	ctx context.Context,
	studentID string, 
	application models.CreatePersonalCollegeApplicationRequestBody,
) (*models.PersonalCollegeApplication, error) {
	createdApplication := &models.PersonalCollegeApplication{}

	const insertQuery = `
	INSERT INTO public.personal_college_applications (
		student_id, global_college_id, application_type, category
	) VALUES (
		$1, $2, $3, $4
	)
	RETURNING id, student_id, global_college_id, application_type, category, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		insertQuery,
		studentID,
		application.GlobalCollegeID,
		application.ApplicationType,
		application.Category,
	).Scan(
		&createdApplication.ID,
		&createdApplication.StudentID,
		&createdApplication.GlobalCollegeID,
		&createdApplication.ApplicationType,
		&createdApplication.Category,
		&createdApplication.CreatedAt,
		&createdApplication.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			
			switch pgErr.Code {
			case pgForeignKeyViolationCode:
				return nil, errs.NotFound("global college", "id", application.GlobalCollegeID)
			case pgCheckViolationCode:
				return nil, errs.BadRequest("invalid application_type or category")
			}
		}
		return nil, err
	}

	return createdApplication, nil
}