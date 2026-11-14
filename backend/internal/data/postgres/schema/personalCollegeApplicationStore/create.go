package personalCollegeApplicationRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"inspirate-consulting/internal/data/postgres/schema"
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
	studentID uuid.UUID,
	application models.CreatePersonalCollegeApplicationRequestBody,
) (*models.PersonalCollegeApplication, error) {
	insertQuery, err := schema.ReadSQLBaseScript("create_personal_college_application.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return nil, err
	}

	var id int64
	err = r.db.QueryRow(
		ctx,
		insertQuery,
		studentID,
		application.GlobalCollegeID,
		application.ApplicationType,
		application.Category,
	).Scan(&id)
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

	return getApplication(ctx, r.db, studentID, id)
}
