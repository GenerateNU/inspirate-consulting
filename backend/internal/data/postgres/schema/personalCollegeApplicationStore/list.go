package personalCollegeApplicationRepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// lists a student's applications (ranked first, in rank order), with cycle_offered and essay_progress computed
func (r *PersonalCollegeApplicationRepository) ListPersonalCollegeApplicationsByStudentID(ctx context.Context, studentID uuid.UUID) ([]models.PersonalCollegeApplication, error) {
	return selectApplications(ctx, r.db, studentID, nil)
}

// gets one of the student's applications, with the same computed fields as the list
func getApplication(ctx context.Context, db dbinterface.QueryInterface, studentID uuid.UUID, applicationID int64) (*models.PersonalCollegeApplication, error) {
	applications, err := selectApplications(ctx, db, studentID, &applicationID)
	if err != nil {
		return nil, err
	}
	if len(applications) == 0 {
		return nil, errs.NotFound("personal college application", "id", applicationID)
	}
	return &applications[0], nil
}

func selectApplications(ctx context.Context, db dbinterface.QueryInterface, studentID uuid.UUID, applicationID *int64) ([]models.PersonalCollegeApplication, error) {
	selectQuery, err := schema.ReadSQLBaseScript("select_personal_college_applications.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, selectQuery, studentID, applicationID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.PersonalCollegeApplication])
}
