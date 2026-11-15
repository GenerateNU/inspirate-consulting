package activity

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// GetActivities returns activity across every student assigned to the
// authenticated counselor, optionally scoped to one of them.
func (h *Handler) GetActivities(ctx context.Context, input *models.GetActivitiesInput) ([]models.Activity, error) {
	counselorID, err := auth.GetCounselorID(ctx)
	if err != nil {
		return nil, err
	}

	var studentID *uuid.UUID
	if input.StudentID != uuid.Nil {
		studentID = &input.StudentID
	}

	return h.ActivityRepository.GetActivitiesForCounselor(ctx, counselorID, studentID)
}
