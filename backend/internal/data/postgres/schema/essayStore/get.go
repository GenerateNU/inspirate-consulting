package essayRepository

import (
	"context"

	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

func (r *EssayRepository) GetEssaysFromStudent(ctx context.Context, studentID uuid.UUID) ([]models.Essays, error) {

	essays := []models.Essays{}

	const selectQuery = ` 
	SELECT id, student_id, type, college_id, link_to_content, status
	FROM public.essays
	WHERE student_id = $1
	`
	rows, err := r.db.Query(ctx, selectQuery, studentID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var essay models.Essays

		err := rows.Scan(
			&essay.ID,
			&essay.StudentID,
			&essay.Type,
			&essay.CollegeID,
			&essay.LinkToContent,
			&essay.Status,
		)

		if err != nil {
			return nil, err
		}

		essays = append(essays, essay)

	}

	// rows.Next returns false both when the rows are exhausted and when the
	// connection fails part way through, so a partial result would otherwise
	// look like a complete one
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return essays, nil

}
