package essayRepository

import (
	"context"

	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

func (r *EssayRepository) UpdateStatus(ctx context.Context, essayID uuid.UUID, status models.Status) error {

	const updateQuery = `
	UPDATE public.essays 
	SET status = $2
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, updateQuery, essayID, status)

	if err != nil {
		return err
	}

	return nil
}
