package mediaRepository

import (
	"context"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *MediaRepository) CreateMedia(ctx context.Context, item *models.CreateMediaRequestBody) (*models.Media, error) {
	createdMedia := &models.Media{}

	insertQuery, err := schema.ReadSQLBaseScript("create_media.sql", SqlMediaFiles)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(
		ctx,
		insertQuery,
		item.Title,
		item.Description,
		item.LengthInMins,
		item.SchoolYear,
		item.S3Key,
	).Scan(
		&createdMedia.ID,
		&createdMedia.Title,
		&createdMedia.Description,
		&createdMedia.LengthInMins,
		&createdMedia.SchoolYear,
		&createdMedia.S3Key,
	)
	if err != nil {
		return nil, err
	}

	return createdMedia, nil
}
