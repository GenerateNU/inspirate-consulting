package mediaAccessRepository

import (
	"context"

	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

func(r *MediaAccessRepository) ListAccessibleMedia(ctx context.Context, studentID string)([]models.Media, error){
	const listQuery = `
	SELECT m.id, m.title, m.description, m.length_in_mins, m.school_year, m.s3_key
	FROM public.media m
	JOIN public.media_access ma ON ma.media_id = m.id
	WHERE ma.student_id = $1
	`
	rows, err := r.db.Query(ctx, listQuery, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accessibleMedia, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Media])

	if err != nil{
		return nil, err
	}

	return accessibleMedia, nil
}