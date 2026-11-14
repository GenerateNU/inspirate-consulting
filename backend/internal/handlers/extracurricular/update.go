package extracurricular

import (
	"context"

	"inspirate-consulting/internal/models"
)

// UpdateExtracurricular updates an extracurricular record and forces student_id from auth context.
func (h *Handler) UpdateExtracurricular(ctx context.Context, id int64, input *models.UpdateExtracurricularInput) (*models.Extracurricular, error) {
	studentID, err := h.ExtracurricularRepository.GetExtracurricularStudentID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := h.Authorizer.CanAccessStudent(ctx, studentID); err != nil {
		return nil, err
	}
	return h.ExtracurricularRepository.UpdateExtracurricular(ctx, id, *input)
}
