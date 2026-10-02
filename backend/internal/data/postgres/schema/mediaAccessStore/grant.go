package mediaAccessRepository

import (
	"context"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *MediaAccessRepository) GrantMediaAccess(ctx context.Context, body *models.GrantMediaAccessRequestBody) (*models.MediaAccess, error) {
	createdAccess := &models.MediaAccess{}

	insertQuery, err := schema.ReadSQLBaseScript("grant_media_access.sql", SqlMediaAccessFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, insertQuery, body.StudentID, body.MediaID).Scan(&createdAccess.ID, &createdAccess.StudentID, &createdAccess.MediaID)

	if err != nil {
		return nil, err
	}

	return createdAccess, nil

}
