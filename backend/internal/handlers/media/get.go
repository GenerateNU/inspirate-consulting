package media

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// returns any media to a counselor, but only media a student has been granted access to.
func (h *Handler) GetMedia(ctx context.Context, id string) (*models.Media, error) {
	if studentID, err := auth.GetStudentID(ctx); err == nil {
		hasAccess, err := h.MediaAccessRepository.HasMediaAccess(ctx, studentID.String(), id)
		if err != nil {
			return nil, err
		}
		if !hasAccess {
			return nil, errs.NotFound("media", "id", id)
		}
	} else if _, err := auth.GetCounselorID(ctx); err != nil {
		return nil, errs.Forbidden("only students and counselors can view media")
	}

	fetchedMedia, err := h.MediaRepository.GetMedia(ctx, id)
	if err != nil {
		return nil, err
	}
	return fetchedMedia, nil
}
