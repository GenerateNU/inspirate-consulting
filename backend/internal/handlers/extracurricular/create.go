package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

// CreateExtracurricular creates a new extracurricular entry using the authenticated student's id from context.
func (h *Handler) CreateExtracurricular(ctx context.Context, input *models.CreateExtracurricularInput) (*models.CreateExtracurricularOutput, error) {
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

	if _, err := uuid.Parse(input.Body.UserID); err != nil {
		input.Body.UserID = studentID
	}

	input.Body.StudentID = &studentID

	if _, err := models.ParseDate(input.Body.StartDate); err != nil {
		return nil, errs.BadRequest(err.Error())
	}

	if input.Body.EndDate != nil {
		if _, err := models.ParseDate(*input.Body.EndDate); err != nil {
			return nil, errs.BadRequest(err.Error())
		}
	}

	extracurricular, err := h.ExtracurricularRepository.CreateExtracurricular(ctx, *input)
	if err != nil {
		return nil, err
	}
	return extracurricular, nil
}
