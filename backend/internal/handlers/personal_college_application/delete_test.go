package personalcollegeapplication

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_DeletePersonalCollegeApplication(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	applicationID := int64(42)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("DeletePersonalCollegeApplication", mock.Anything, mock.Anything, applicationID).Return(nil)

		handler := NewHandler(mockPCARepo, nil)
		err := handler.DeletePersonalCollegeApplication(ctx, applicationID)

		assert.NoError(t, err)
	})

	t.Run("repository error propagates (e.g. not found)", func(t *testing.T) {
		t.Parallel()

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("DeletePersonalCollegeApplication", mock.Anything, mock.Anything, applicationID).
			Return(errors.New("personal college application not found"))

		handler := NewHandler(mockPCARepo, nil)
		err := handler.DeletePersonalCollegeApplication(ctx, applicationID)

		assert.Error(t, err)
		assert.EqualError(t, err, "personal college application not found")
	})
}
