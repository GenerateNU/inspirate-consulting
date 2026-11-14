package essays

import (
	"context"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) UpdateEssayGroup(ctx context.Context, input *models.UpdateEssayGroupInput) (*models.UpdateEssayGroupOutput, error) {
	studentID, err := h.EssayRepository.GetEssayStudentID(ctx, input.EssayID)
	if err != nil {
		return nil, err
	}

	if err := h.Authorizer.CanAccessStudent(ctx, studentID); err != nil {
		return nil, err
	}

	if input.Body.EssayGroupID != nil {
		groupStudentID, err := h.EssayGroupsRepository.GetEssayGroupStudentID(ctx, *input.Body.EssayGroupID)
		if err != nil {
			return nil, err
		}
		if groupStudentID != studentID {
			return nil, errs.Forbidden("essay group does not belong to this student")
		}
	}

	updatedEssay, err := h.EssayRepository.UpdateEssayGroup(ctx, input.EssayID, input.Body.EssayGroupID)
	if err != nil {
		return nil, err
	}

	return &models.UpdateEssayGroupOutput{
		Body: models.EssayBody{Essay: updatedEssay},
	}, nil
}
