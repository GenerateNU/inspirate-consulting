package personalCollegeApplicationRepository

import (
	"context"
	"errors"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// updates an application's college, cycle, category, interest, and status (School Details save)
func (r *PersonalCollegeApplicationRepository) UpdatePersonalCollegeApplication(
	ctx context.Context,
	studentID uuid.UUID,
	applicationID int64,
	input models.UpdatePersonalCollegeApplicationRequestBody,
) (*models.PersonalCollegeApplication, error) {
	updateQuery, err := schema.ReadSQLBaseScript("update_personal_college_application_general.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return nil, err
	}

	var id int64
	err = r.db.QueryRow(
		ctx,
		updateQuery,
		applicationID,
		studentID,
		input.GlobalCollegeID,
		input.ApplicationType,
		input.Category,
		input.Interest,
		input.ApplicationStatus,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("personal college application", "id", applicationID)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgForeignKeyViolationCode:
				return nil, errs.NotFound("global college", "id", input.GlobalCollegeID)
			case pgCheckViolationCode:
				return nil, errs.BadRequest("invalid application_type, category, interest, or application_status")
			}
		}
		return nil, err
	}

	return getApplication(ctx, r.db, studentID, id)
}
