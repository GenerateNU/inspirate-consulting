package essay_groups

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateEssayGroup(ctx context.Context, input *models.CreateEssayGroupInput) (*models.CreateEssayGroupOutput, error) {
	studentID, err := h.Authorizer.ResolveStudentID(ctx, input.Body.StudentID)
	if err != nil {
		return nil, err
	}
	createdGroup, err := h.EssayGroupRepository.CreateEssayGroup(ctx, studentID, input.Body)
	if err != nil {
		return nil, err
	}

	return &models.CreateEssayGroupOutput{
		Body: *createdGroup,
	}, nil
}
