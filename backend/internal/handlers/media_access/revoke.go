package mediaaccess

import (
	"context"
)

func (h *Handler) RevokeMediaAccess(ctx context.Context, id string) error {
	studentID, err := h.MediaAccessRepository.GetMediaAccessStudentID(ctx, id)
	if err != nil {
		return err
	}

	if err := h.Authorizer.CanAccessStudent(ctx, studentID); err != nil {
		return err
	}

	return h.MediaAccessRepository.RevokeMediaAccess(ctx, id)
}
