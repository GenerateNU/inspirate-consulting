package mediaaccess

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the RevokeMediaAccess handler logic without HTTP or database dependencies.
func TestHandler_RevokeMediaAccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RevokeMediaAccess", mock.Anything, id).Return(nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		err := handler.RevokeMediaAccess(ctx, id)

		// Verify result
		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RevokeMediaAccess", mock.Anything, id).Return(errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		err := handler.RevokeMediaAccess(ctx, id)

		// Verify error propagation
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}
