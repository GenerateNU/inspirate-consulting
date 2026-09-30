package mediaRepository

import (
	"context"
	"errors"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *MediaRepository) GetMedia(ctx context.Context, id string) (*models.Media, error) {
	getMedia := &models.Media{}

	selectQuery, err := schema.ReadSQLBaseScript("get_media.sql", SqlMediaFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, selectQuery, id).Scan(
		&getMedia.ID,
		&getMedia.Title,
		&getMedia.Description,
		&getMedia.LengthInMins,
		&getMedia.SchoolYear,
		&getMedia.S3Key,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("media", "id", id)
		}

		return nil, err
	}

	return getMedia, nil

}
