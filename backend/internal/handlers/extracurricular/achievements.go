package extracurricular

import (
	"context"

	"inspirate-consulting/internal/models"
)

// CreateExtracurricularAchievement adds an achievement to an extracurricular the caller can access.
func (h *Handler) CreateExtracurricularAchievement(ctx context.Context, input *models.CreateExtracurricularAchievementInput) (*models.ExtracurricularAchievement, error) {
	if _, err := h.authorizeExtracurricular(ctx, input.ExtracurricularID); err != nil {
		return nil, err
	}
	return h.ExtracurricularRepository.CreateExtracurricularAchievement(ctx, input.ExtracurricularID, input.Body)
}

// UpdateExtracurricularAchievement edits an achievement on an extracurricular the caller can access.
func (h *Handler) UpdateExtracurricularAchievement(ctx context.Context, input *models.UpdateExtracurricularAchievementInput) (*models.ExtracurricularAchievement, error) {
	if _, err := h.authorizeExtracurricular(ctx, input.ExtracurricularID); err != nil {
		return nil, err
	}
	return h.ExtracurricularRepository.UpdateExtracurricularAchievement(ctx, input.ExtracurricularID, input.AchievementID, input.Body)
}

// DeleteExtracurricularAchievement removes an achievement from an extracurricular the caller can access.
func (h *Handler) DeleteExtracurricularAchievement(ctx context.Context, input *models.DeleteExtracurricularAchievementInput) error {
	if _, err := h.authorizeExtracurricular(ctx, input.ExtracurricularID); err != nil {
		return err
	}
	return h.ExtracurricularRepository.DeleteExtracurricularAchievement(ctx, input.ExtracurricularID, input.AchievementID)
}
