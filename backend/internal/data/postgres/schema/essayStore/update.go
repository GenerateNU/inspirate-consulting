package essayRepository

import (
	"context"

	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *EssayRepository) UpdateStatus(ctx context.Context, essayID uuid.UUID, status models.Status) (*models.Essays, error) {

	const updateQuery = `
	UPDATE public.essays 
	SET status = $2
	WHERE id = $1
	RETURNING id, student_id, type, college_id, link_to_content, status
	`

	rows, err := r.db.Query(ctx, updateQuery, essayID, status)
	if err != nil {
		return nil, err
	}

	essay, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[models.Essays])
	if err != nil {
		return nil, err
	}

	return &essay, nil
}
