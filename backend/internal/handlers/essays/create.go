package essays

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateEssay(ctx context.Context, input *models.CreateEssayInput) (*models.CreateEssayOutput, error) {
	essay := models.Essays{
		StudentID:     input.Body.StudentID,
		Type:          input.Body.Type,
		CollegeID:     input.Body.CollegeID,
		LinkToContent: input.Body.LinkToContent,
	}

	if err := h.EssayRepository.CreateEssay(ctx, essay); err != nil {
		return nil, err
	}

	return &models.CreateEssayOutput{}, nil
}
