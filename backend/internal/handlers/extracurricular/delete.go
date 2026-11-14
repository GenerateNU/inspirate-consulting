package extracurricular

import (
	"context"
)

// DeleteExtracurricular deletes an extracurricular and, by cascade, its
// achievements and notes. Students may only delete their own; counselors only
// those of assigned students. The owner is read from the existing row.
func (h *Handler) DeleteExtracurricular(ctx context.Context, id int64) error {
	existing, err := h.authorizeExtracurricular(ctx, id)
	if err != nil {
		return err
	}

	return h.ExtracurricularRepository.DeleteExtracurricular(ctx, id, existing.StudentID)
}
