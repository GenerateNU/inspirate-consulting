package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// authorizeExtracurricular loads an extracurricular and checks that the caller
// may access its owner: students only their own, counselors only assigned
// students. Returns a 404 HTTPError if the extracurricular does not exist.
func (h *Handler) authorizeExtracurricular(ctx context.Context, id int64) (*models.ExtracurricularDetail, error) {
	existing, err := h.ExtracurricularRepository.GetExtracurricular(ctx, id)
	if err != nil {
		return nil, err
	}

	ownerID, err := auth.ParseStudentID(existing.StudentID)
	if err != nil {
		return nil, err
	}
	if err := h.Authorizer.CanAccessStudent(ctx, ownerID); err != nil {
		return nil, err
	}

	return existing, nil
}
