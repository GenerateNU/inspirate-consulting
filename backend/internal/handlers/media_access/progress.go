package mediaaccess

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// marks media as started (when the student opens it) or completed (when they mark
// it done) for the calling student. Idempotent: repeating a status keeps the first timestamp.
func (h *Handler) RecordMediaProgress(ctx context.Context, mediaID string, status models.MediaProgressStatus) (*models.MediaAccess, error) {
	studentID, err := auth.GetStudentID(ctx)
	if err != nil {
		return nil, errs.Forbidden("only students can record media progress")
	}

	var completed bool
	switch status {
	case models.MediaProgressStarted:
		completed = false
	case models.MediaProgressCompleted:
		completed = true
	default:
		return nil, errs.BadRequest("status must be started or completed")
	}

	access, err := h.MediaAccessRepository.RecordMediaProgress(ctx, studentID.String(), mediaID, completed)
	if err != nil {
		return nil, err
	}
	return access, nil
}
