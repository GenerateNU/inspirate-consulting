package mediaaccess

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// minimum pg_trgm word_similarity, from 0.0 (complete mismatch) to 1.0 (perfect match), between the
// search query and the best-matching part of a media's title or description
const searchThreshold = 0.4

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
