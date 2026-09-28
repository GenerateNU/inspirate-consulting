package mediaRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *MediaRepository) GetMedia(ctx context.Context, id string) (*models.Media, error){
	getMedia := &models.Media{}

	const selectQuery = `
	SELECT id, title, description, length_in_mins, school_year, s3_key
	FROM public.key
	WHERE id = $1
	`

	err := r.db.QueryRow(ctx, selectQuery, id).Scan(
		&getMedia.ID,
		&getMedia.Title,
		&getMedia.Description,
		&getMedia.LengthInMins,
		&getMedia.SchoolYear,
		&getMedia.S3Key,
	)
	if(err != nil){
		if errors.Is(err, pgx.ErrNoRows){
			return nil, errs.NotFound("media", "id", id)
		}

		return nil, err
	}

	return getMedia, nil
	
}