package mediaAccessRepository

import (
	"context"
	"errors"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

// marks the student's access row as started (and completed, if completed is true)
// Repeat calls keep the original timestamps
func (r *MediaAccessRepository) RecordMediaProgress(ctx context.Context, studentID string, mediaID string, completed bool) (*models.MediaAccess, error) {
	query, err := schema.ReadSQLBaseScript("record_media_progress.sql", SqlMediaAccessFiles)
	if err != nil {
		return nil, err
	}

	access := &models.MediaAccess{}
	err = r.db.QueryRow(ctx, query, studentID, mediaID, completed).Scan(
		&access.ID,
		&access.StudentID,
		&access.MediaID,
		&access.StartedAt,
		&access.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("media", "id", mediaID)
		}
		return nil, err
	}
	return access, nil
}
