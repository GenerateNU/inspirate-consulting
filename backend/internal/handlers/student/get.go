package student

import (
	"context"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

func (h *Handler) GetStudent(ctx context.Context, id uuid.UUID) (*models.Student, error) {
	return h.StudentRepository.GetStudent(ctx, id)
}
