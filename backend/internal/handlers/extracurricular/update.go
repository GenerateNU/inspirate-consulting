package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// UpdateExtracurricular updates an extracurricular record and forces student_id from auth context.
func (h *Handler) UpdateExtracurricular(ctx context.Context, id int64, input *models.UpdateExtracurricularInput) (*models.Extracurricular, error) {
	if auth.GetStudentID(ctx) == "" {
		return nil, errs.BadRequest("student id not found in request context")
	}

	return h.ExtracurricularRepository.UpdateExtracurricular(ctx, id, *input)
}