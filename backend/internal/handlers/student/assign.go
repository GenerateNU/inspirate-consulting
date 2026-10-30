package student

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) AssignCounselor(ctx context.Context, input *models.AssignCounselorInput) (*models.Student, error) {

	return h.StudentRepository.AssignCounselor(ctx, input.Body.UserID, input.Body.CounselorID)
}
