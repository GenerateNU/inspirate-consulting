package essays

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) GetEssaysFromStudent(ctx context.Context, input *models.GetEssaysFromStudentInput) (*models.GetEssaysFromStudentOutput, error) {
	essays, err := h.EssayRepository.GetEssaysFromStudent(ctx, input.StudentID)
	if err != nil {
		return nil, err
	}

	return &models.GetEssaysFromStudentOutput{
		Body: models.EssayListBody{Essays: essays},
	}, nil
}
