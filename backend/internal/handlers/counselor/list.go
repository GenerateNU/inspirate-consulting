package counselor

import (
	"context"

	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"
)

func (h *Handler) ListCounselors(ctx context.Context, input *models.FetchCounselorsInput) (*models.FetchCounselorsBody, error) {

	var before *pagination.NameIDKey
	if input.Cursor != "" {
		key, err := pagination.Decode[pagination.NameIDKey](input.Cursor)
		if err != nil {
			return nil, err
		}
		before = &key
	}

	rows, err := h.CounselorRepository.ListCounselors(ctx, &input.Name, before, input.Limit+1)
	if err != nil {
		return nil, err
	}

	counselors, nextCursor, err := pagination.NextPage(rows, input.Limit, func(input models.CounselorListItem) pagination.NameIDKey {
		return pagination.NameIDKey{Name: input.Name, ID: input.ID}
	})
	if err != nil {
		return nil, err
	}

	return &models.FetchCounselorsBody{Counselors: counselors,
		NextCursor: nextCursor}, nil
}
