package personalCollegeApplicationRepository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *PersonalCollegeApplicationRepository) ListPersonalCollegeApplicationsByStudentID(ctx context.Context, studentID string) ([]models.PersonalCollegeApplication, error) {
	listQuery, err := schema.ReadSQLBaseScript("list_personal_college_applications_by_student.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, listQuery, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applications, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.PersonalCollegeApplication])
	if err != nil {
		return nil, err
	}

	return applications, nil
}
