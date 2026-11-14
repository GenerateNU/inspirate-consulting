package essays

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) UpdateStatus(ctx context.Context, input *models.UpdateStatusInput) (*models.UpdateStatusOutput, error) {
	studentID, err := h.EssayRepository.GetEssayStudentID(ctx, input.EssayID)

	if err != nil {
		return nil, err
	}

	if err := h.Authorizer.CanAccessStudent(ctx, studentID); err != nil {
		return nil, err
	}

	updatedEssay, err := h.EssayRepository.UpdateStatus(ctx, input.EssayID, input.Body.Status)
	if err != nil {
		return nil, err
	}

	return &models.UpdateStatusOutput{
		Body: models.EssayBody{Essay: updatedEssay},
	}, nil

}
