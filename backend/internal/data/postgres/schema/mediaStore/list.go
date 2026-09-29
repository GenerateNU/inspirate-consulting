package mediaRepository

import(
	"context"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/models"
)

func(r *MediaRepository) ListAllMedia(ctx context.Context)([]models.Media, error){
	const listQuery = `
	SELECT id, title, description, length_in_mins, school_year, s3_key
	FROM public.media
	`

	rows, err := r.db.Query(ctx, listQuery)

	if(err != nil){
		return nil, err
	}
	defer rows.Close()

	allMedia, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Media])
	
	if(err != nil){
		return nil, err
	}
	return allMedia, nil
}