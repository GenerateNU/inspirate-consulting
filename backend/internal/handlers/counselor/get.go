package counselor

import (
	"context"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

// fetches a counselor's profile; like GET /students/{id}, any logged-in user can read it
func (h *Handler) GetCounselor(ctx context.Context, id uuid.UUID) (*models.Counselor, error) {
	return h.CounselorRepository.GetCounselor(ctx, id)
}
