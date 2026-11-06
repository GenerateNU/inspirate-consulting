package mediaRepository

import (
	"context"
	"errors"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *MediaRepository) UpdateMedia(ctx context.Context, id string, change string) (*models.Media, error) {
	query, err := schema.ReadSQLBaseScript("update_media.sql", SqlMediaFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	updated := &models.Media{}
	err = r.db.QueryRow(ctx, query, id, change).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Description,
		&updated.LengthInMins,
		&updated.SchoolYear,
		&updated.S3Key,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("media", "id", id)
		}
		return nil, err
	}

	return updated, nil

}
