package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// GetExtracurricular returns one extracurricular with its achievements and notes.
// Students may only read their own; counselors only those of assigned students.
func (h *Handler) GetExtracurricular(ctx context.Context, id int64) (*models.ExtracurricularDetail, error) {
	detail, err := h.ExtracurricularRepository.GetExtracurricular(ctx, id)
	if err != nil {
		return nil, err
	}

	ownerID, err := auth.ParseStudentID(detail.StudentID)
	if err != nil {
		return nil, err
	}
	if err := h.Authorizer.CanAccessStudent(ctx, ownerID); err != nil {
		return nil, err
	}

	return detail, nil
}
