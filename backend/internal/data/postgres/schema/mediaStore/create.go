package mediaRepository
import (
	"context"
	"inspirate-consulting/internal/models"
)

func(r *MediaRepository) CreateMedia(ctx context.Context, item *models.CreateMediaRequestBody)(*models.Media, error){
	createdMedia := &models.Media{}

	const insertQuery = `
	INSERT INTO public.media (
		title, description, length_in_mins, school_year, s3_key
	) VALUES(
	 	$1, $2, $3, $4, $5
	)
	RETURNING id, title, description, length_in_mins, school_year, s3_key
	`
	err := r.db.QueryRow(
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
