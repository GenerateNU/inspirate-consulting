package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
)

// DeleteExtracurricular deletes an extracurricular and, by cascade, its
// achievements and notes. Students may only delete their own; counselors only
// those of assigned students. The owner is read from the existing row.
func (h *Handler) DeleteExtracurricular(ctx context.Context, id int64) error {
	existing, err := h.ExtracurricularRepository.GetExtracurricular(ctx, id)
	if err != nil {
		return err
	}

	ownerID, err := auth.ParseStudentID(existing.StudentID)
	if err != nil {
		return err
	}
	if err := h.Authorizer.CanAccessStudent(ctx, ownerID); err != nil {
		return err
	}

	return h.ExtracurricularRepository.DeleteExtracurricular(ctx, id, existing.StudentID)
}
