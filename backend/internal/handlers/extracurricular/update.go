package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// UpdateExtracurricular updates an extracurricular record and forces student_id from auth context.
func (h *Handler) UpdateExtracurricular(ctx context.Context, id int64, input *models.UpdateExtracurricularInput) (*models.Extracurricular, error) {
	if _, err := auth.GetStudentID(ctx); err != nil {
		return nil, errs.Unauthorized()
	}

	return h.ExtracurricularRepository.UpdateExtracurricular(ctx, id, *input)
}
