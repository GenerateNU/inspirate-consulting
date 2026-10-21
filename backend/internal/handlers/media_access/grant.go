package mediaaccess

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

func (h *Handler) GrantMediaAccess(ctx context.Context, body *models.GrantMediaAccessRequestBody) (*models.MediaAccess, error) {
	studentID, err := auth.ParseStudentID(body.StudentID)
	if err != nil {
		return nil, err
	}
	if err := h.Authorizer.CanAccessStudent(ctx, studentID); err != nil {
		return nil, err
	}

	createdAccess, err := h.MediaAccessRepository.GrantMediaAccess(ctx, body)
	if err != nil {
		return nil, err
	}

	return createdAccess, nil
}
