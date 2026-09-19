package globalCollegeRepository

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (r *GlobalCollegeRepository) ListGlobalColleges(ctx context.Context) (*models.ListGlobalCollegesOutput, error) {

	const listQuery = `
	SELECT id, created_at, updated_at, school_name, school_location, ea_deadline, ed_deadline, rd_deadline
	FROM public.global_colleges
	`

	rows, err := r.db.Query(ctx, listQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	globalColleges := []models.GlobalCollege{}

	for rows.Next() {
		var globalCollege models.GlobalCollege
		if err := rows.Scan(
			&globalCollege.ID,
			&globalCollege.CreatedAt,
			&globalCollege.UpdatedAt,
			&globalCollege.SchoolName,
			&globalCollege.SchoolLocation,
			&globalCollege.EADeadline,
			&globalCollege.EDDeadline,
			&globalCollege.RDDeadline,
		); err != nil {
			return nil, err
		}
		globalColleges = append(globalColleges, globalCollege)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.ListGlobalCollegesOutput{Body: globalColleges}, nil
}