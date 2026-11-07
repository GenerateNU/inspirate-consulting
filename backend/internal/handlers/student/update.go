package student

import (
	"context"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

// updates a student's profile fields
// the student themselves or their assigned counselor can do this
func (h *Handler) UpdateStudent(ctx context.Context, id uuid.UUID, input models.UpdateStudentBody) (*models.Student, error) {
	if err := h.Authorizer.CanAccessStudent(ctx, id); err != nil {
		return nil, err
	}

	updatedStudent, err := h.StudentRepository.UpdateStudent(ctx, id, input)
	if err != nil {
		return nil, err
	}

	return updatedStudent, nil
}
