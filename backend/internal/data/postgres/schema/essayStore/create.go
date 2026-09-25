package essayRepository

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (r *EssayRepository) CreateEssay(ctx context.Context, essay models.Essays) error {
	const insertQuery = `
	INSERT INTO public.essays (student_id, type, college_id, link_to_content)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, insertQuery, essay.StudentID, essay.Type, essay.CollegeID, essay.LinkToContent)

	if err != nil {
		return err
	}

	return nil
}
