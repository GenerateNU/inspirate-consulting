package extracurricular

import (
	"context"

	"inspirate-consulting/internal/models"
)

// CreateExtracurricularNote adds a note to an extracurricular the caller can access.
func (h *Handler) CreateExtracurricularNote(ctx context.Context, input *models.CreateExtracurricularNoteInput) (*models.ExtracurricularNote, error) {
	if _, err := h.authorizeExtracurricular(ctx, input.ExtracurricularID); err != nil {
		return nil, err
	}
	return h.ExtracurricularRepository.CreateExtracurricularNote(ctx, input.ExtracurricularID, input.Body)
}

// UpdateExtracurricularNote edits a note on an extracurricular the caller can access.
func (h *Handler) UpdateExtracurricularNote(ctx context.Context, input *models.UpdateExtracurricularNoteInput) (*models.ExtracurricularNote, error) {
	if _, err := h.authorizeExtracurricular(ctx, input.ExtracurricularID); err != nil {
		return nil, err
	}
	return h.ExtracurricularRepository.UpdateExtracurricularNote(ctx, input.ExtracurricularID, input.NoteID, input.Body)
}

// DeleteExtracurricularNote removes a note from an extracurricular the caller can access.
func (h *Handler) DeleteExtracurricularNote(ctx context.Context, input *models.DeleteExtracurricularNoteInput) error {
	if _, err := h.authorizeExtracurricular(ctx, input.ExtracurricularID); err != nil {
		return err
	}
	return h.ExtracurricularRepository.DeleteExtracurricularNote(ctx, input.ExtracurricularID, input.NoteID)
}
