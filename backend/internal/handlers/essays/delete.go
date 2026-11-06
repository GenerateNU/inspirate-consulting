package essays

import (
	"context"

	"github.com/google/uuid"
)

// DeleteEssay lets a student delete their own essay, or a counselor delete an
// essay belonging to one of their assigned students.
func (h *Handler) DeleteEssay(ctx context.Context, essayID uuid.UUID) error {
	essay, err := h.EssayRepository.GetEssay(ctx, essayID)
	if err != nil {
		return err
	}
	if err := h.Authorizer.CanAccessStudent(ctx, essay.StudentID); err != nil {
		return err
	}

	return h.EssayRepository.DeleteEssay(ctx, essayID, essay.StudentID)
}
