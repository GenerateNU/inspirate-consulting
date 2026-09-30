package essayRepository

import (
	"context"
	"fmt"

	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *EssayRepository) GetEssaysFromStudent(ctx context.Context, studentID uuid.UUID) ([]models.Essays, error) {
	const selectQuery = ` 
	SELECT id, student_id, type, college_id, link_to_content, status
	FROM public.essays
	WHERE student_id = $1
	`
	rows, err := r.db.Query(ctx, selectQuery, studentID)
	if err != nil {
		return nil, err
	}

	essays, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Essays])

	if err != nil {
		return nil, fmt.Errorf("failed to collect essays: %w", err)
	}

	return essays, nil
}
