package extracurricular

import (
	"context"
	"log/slog"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// DeleteExtracurricular deletes an extracurricular and, by cascade, its
// achievements and notes. Students may only delete their own; counselors only
// those of assigned students. The owner is read from the existing row.
func (h *Handler) DeleteExtracurricular(ctx context.Context, id int64) error {
	existing, err := h.authorizeExtracurricular(ctx, id)
	if err != nil {
		return err
	}

	if err := h.ExtracurricularRepository.DeleteExtracurricular(ctx, id, existing.StudentID); err != nil {
		return err
	}

	if ownerID, parseErr := auth.ParseStudentID(existing.StudentID); parseErr == nil {
		if err := h.ActivityRepository.CreateActivity(ctx, ownerID, models.CreateActivityBody{
			Action:            models.ActionRemovedExtracurricular,
			ExtracurricularID: &id,
		}); err != nil {
			slog.Error("failed to log removed_extracurricular activity", "student_id", ownerID, "extracurricular_id", id, "error", err)
		}
	} else {
		slog.Error("failed to parse student id for removed_extracurricular activity", "raw", existing.StudentID, "error", parseErr)
	}

	return nil
}
