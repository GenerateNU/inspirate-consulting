package mediaaccess

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the GrantMediaAccess handler logic without HTTP or database dependencies.
func TestHandler_GrantMediaAccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	studentID := uuid.NewString()
	mediaID := uuid.NewString()

	input := &models.GrantMediaAccessRequestBody{
		StudentID: studentID,
		MediaID:   mediaID,
	}
	expectedOutput := &models.MediaAccess{
		ID:        uuid.NewString(),
		StudentID: studentID,
		MediaID:   mediaID,
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GrantMediaAccess", mock.Anything, input).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.GrantMediaAccess(ctx, input)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GrantMediaAccess", mock.Anything, input).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.GrantMediaAccess(ctx, input)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
