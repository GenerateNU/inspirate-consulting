package mediaAccessRepository 

import(
	"context"
	"inspirate-consulting/internal/models"
)

func(r *MediaAccessRepository) GrantMediaAccess(ctx context.Context, body *models.GrantMediaAccessRequestBody)(*models.MediaAccess,error){
	createdAccess := &models.MediaAccess{}

	const insertQuery = `
	INSERT INTO public.media_access(
		student_id, media_id
	) VALUES(
		$1, $2
	)
	RETURNING id, student_id, media_id
	`

	err := r.db.QueryRow(ctx, insertQuery, body.StudentID, body.MediaID,).Scan(&createdAccess.ID, &createdAccess.StudentID, &createdAccess.MediaID,)

	if(err != nil){
		return nil, err
	}

	return createdAccess, nil

}