package essays

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) GetEssaysFromStudent(ctx context.Context, input *models.GetEssaysFromStudentInput) (*models.GetEssaysFromStudentOutput, error) {
	if err := h.Authorizer.CanAccessStudent(ctx, input.StudentID); err != nil {
		return nil, err
	}

	filter := models.GetEssaysFromStudentFilter{}
	if input.CollegeID != 0 {
		filter.CollegeID = &input.CollegeID
	}

	essays, err := h.EssayRepository.GetEssaysFromStudent(ctx, input.StudentID, filter)
	if err != nil {
		return nil, err
	}

	return &models.GetEssaysFromStudentOutput{
		Body: models.EssayListBody{Essays: essays},
	}, nil
}
