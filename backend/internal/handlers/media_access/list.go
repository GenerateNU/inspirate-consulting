package mediaaccess

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// minimum pg_trgm similarity, from 0.0 (complete mismatch) to 1.0 (perfect match), for a media's
// title or description to count as matching the search query
const searchThreshold = 0.2

func (h *Handler) ListAccessibleMedia(ctx context.Context, query string, limit int, offset int) ([]models.AccessibleMedia, error) {
	studentID, err := auth.GetStudentID(ctx)
	if err != nil {
		return nil, errs.Forbidden("only students can list their accessible media")
	}

	accessibleMedia, err := h.MediaAccessRepository.ListAccessibleMedia(ctx, studentID.String(), query, searchThreshold, limit, offset)
	if err != nil {
		return nil, err
	}
	return accessibleMedia, nil
}
