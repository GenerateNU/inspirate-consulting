package essay_groups

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) GetEssayGroups(ctx context.Context, input *models.GetEssayGroupsInput) (*models.GetEssayGroupsOutput, error) {
	if err := h.Authorizer.CanAccessStudent(ctx, input.StudentID); err != nil {
		return nil, err
	}

	essayGroups, err := h.EssayGroupRepository.ListEssayGroups(ctx, input.StudentID)
	if err != nil {
		return nil, err
	}

	return &models.GetEssayGroupsOutput{
		Body: models.EssayGroupListBody{EssayGroups: essayGroups},
	}, nil
}
