package essays

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateEssay(ctx context.Context, input *models.CreateEssayInput) (*models.CreateEssayOutput, error) {
	studentID, err := auth.GetStudentID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.EssayRepository.CreateEssay(ctx, studentID, input.Body); err != nil {
		return nil, err
	}

	return &models.CreateEssayOutput{}, nil
}
