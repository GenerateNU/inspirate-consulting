package mediaaccess

import (
	"context"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/auth"
)

func (h *Handler) ListAccessibleMedia(ctx context.Context, limit int, offset int)([]models.Media, error){
	studentID := auth.GetStudentID(ctx)

	accessibleMedia, err := h.MediaAccessRepository.ListAccessibleMedia(ctx, studentID, limit, offset)
	if(err != nil){
		return nil, err
	}
	return accessibleMedia, nil

}
