package personalCollegeApplicationRepository

import (
	"context"
	"errors"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

// updates an application's global_college_id, application_type, and category
func (r *PersonalCollegeApplicationRepository) UpdatePersonalCollegeApplication(
	ctx context.Context,
	studentID string,
	applicationID int64,
	input models.UpdatePersonalCollegeApplicationRequestBody,
) (*models.PersonalCollegeApplication, error) {
	updatedApplication := &models.PersonalCollegeApplication{}

	updateQuery, err := schema.ReadSQLBaseScript("update_personal_college_application_general.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		updateQuery,
		applicationID,
		studentID,
		input.GlobalCollegeID,
		input.ApplicationType,
		input.Category,
	).Scan(
		&updatedApplication.ID,
		&updatedApplication.StudentID,
		&updatedApplication.GlobalCollegeID,
		&updatedApplication.ApplicationType,
		&updatedApplication.Category,
		&updatedApplication.Rank,
		&updatedApplication.CreatedAt,
		&updatedApplication.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("personal college application", "id", applicationID)
		}
		return nil, err
	}

	return updatedApplication, nil
}
