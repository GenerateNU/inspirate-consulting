package personalCollegeApplicationRepository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/models"
)

func (r *PersonalCollegeApplicationRepository) ListPersonalCollegeApplicationsByStudentID(ctx context.Context, studentID string,) ([]models.PersonalCollegeApplication, error) {
	const listQuery = `
	SELECT id, student_id, global_college_id, application_type, category, created_at, updated_at
	FROM public.personal_college_applications
	WHERE student_id = $1
	`

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