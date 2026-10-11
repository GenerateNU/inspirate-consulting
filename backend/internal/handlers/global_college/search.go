package globalcollege

import (
	"context"

	"inspirate-consulting/internal/models"
)

// range from 0.0 to 1.0, inclusive, where 0.0 is a complete mismatch and 1.0 is a perfect match
const searchThreshold = 0.2

// finds colleges whose name OR location loosely matches the given query, ordered most-relevant first.
func (h *Handler) SearchGlobalColleges(ctx context.Context, query string) ([]models.GlobalCollege, error) {
	results, err := h.GlobalCollegeRepository.SearchGlobalColleges(ctx, query, searchThreshold)
	if err != nil {
		return nil, err
	}

	return results, nil
}
