package mediaaccess

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

func (h *Handler) ListAccessibleMedia(ctx context.Context, limit int, offset int) ([]models.Media, error) {
	studentID, err := auth.GetStudentID(ctx)
	if err != nil {
		return nil, err
	}

	accessibleMedia, err := h.MediaAccessRepository.ListAccessibleMedia(ctx, studentID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	return accessibleMedia, nil

}
