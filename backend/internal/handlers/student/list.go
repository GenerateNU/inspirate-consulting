package student

import (
	"context"

	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"
)

func (h *Handler) ListStudents(ctx context.Context, input *models.FetchStudentsInput) (*models.FetchStudentsBody, error) {
	var before *pagination.NameIDKey
	if input.Cursor != "" {
		key, err := pagination.Decode[pagination.NameIDKey](input.Cursor)
		if err != nil {
			return nil, err
		}
		before = &key
	}

	rows, err := h.StudentRepository.ListStudents(ctx, &input.Name, before, input.Limit+1)
	if err != nil {
		return nil, err
	}

	students, nextCursor, err := pagination.NextPage(rows, input.Limit, func(input models.StudentListItem) pagination.NameIDKey {
		return pagination.NameIDKey{Name: input.Name, ID: input.ID}
	})
	if err != nil {
		return nil, err
	}

	return &models.FetchStudentsBody{Students: students,
		NextCursor: nextCursor}, nil
}
