package essayreview

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// Unit tests for the GetEssayReviewStatus handler logic without HTTP or database dependencies.
func TestHandler_GetEssayReviewStatus(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("success returns the status for the essay", func(t *testing.T) {
		t.Parallel()

		essayID := uuid.New()
		expectedEntity := &models.EssayReviewStatus{
			TransactionID: uuid.New(),
			EssayID:       essayID,
			StudentID:     uuid.New(),
			Status:        models.ReviewStatusOpen,
			RequestedAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockRepo := mocks.NewEssayReviewRepository(t)
		mockRepo.On("DB").Return(nil)
		mockRepo.On("FindEssayReviewStatus", mock.Anything, mock.Anything, essayID).Return(expectedEntity, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.GetEssayReviewStatus(ctx, essayID)

		assert.NoError(t, err)
		assert.Equal(t, expectedEntity, res)
	})

	t.Run("unreviewed essay propagates the repository error", func(t *testing.T) {
		t.Parallel()

		essayID := uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		mockRepo.On("DB").Return(nil)
		mockRepo.On("FindEssayReviewStatus", mock.Anything, mock.Anything, essayID).
			Return(nil, errs.NotFound("essay review", "essay_id", essayID.String()))

		handler := NewHandler(mockRepo)
		res, err := handler.GetEssayReviewStatus(ctx, essayID)

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 404, err.(errs.HTTPError).Code)
	})
}
