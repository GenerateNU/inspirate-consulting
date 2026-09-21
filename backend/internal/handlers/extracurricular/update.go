package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// UpdateExtracurricular updates an extracurricular record and forces student_id from auth context.
func (h *Handler) UpdateExtracurricular(ctx context.Context, id int64, input *models.UpdateExtracurricularInput) (*models.Extracurricular, error) {
	if input == nil {
		return nil, errs.BadRequest("missing request body")
	}

	studentID := auth.GetStudentID(ctx)
	if studentID == "" && input.Body.StudentID != nil && *input.Body.StudentID != "" {
		studentID = *input.Body.StudentID
	}
	if studentID == "" {
		return nil, errs.BadRequest("student id not found in request context")
	}

	input.Body.StudentID = &studentID
	if input.Body.StartDate != nil {
		if _, err := models.ParseDate(*input.Body.StartDate); err != nil {
			return nil, errs.BadRequest(err.Error())
		}
	}
	if input.Body.EndDate != nil {
		if _, err := models.ParseDate(*input.Body.EndDate); err != nil {
			return nil, errs.BadRequest(err.Error())
		}
	}

	input.ID = id
	return h.ExtracurricularRepository.UpdateExtracurricular(ctx, id, *input)
}
