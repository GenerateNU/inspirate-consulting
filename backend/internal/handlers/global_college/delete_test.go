package globalcollege

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_DeleteGlobalCollege(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	collegeID := int64(1)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("DeleteGlobalCollege", mock.Anything, collegeID).Return(nil)

		handler := NewHandler(mockRepo)
		err := handler.DeleteGlobalCollege(ctx, collegeID)

		assert.NoError(t, err)
	})

	t.Run("repository error propagates (e.g. not found)", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("DeleteGlobalCollege", mock.Anything, collegeID).
			Return(errors.New("global college not found"))

		handler := NewHandler(mockRepo)
		err := handler.DeleteGlobalCollege(ctx, collegeID)

		assert.Error(t, err)
		assert.EqualError(t, err, "global college not found")
	})

	t.Run("repository error propagates (e.g. still referenced, conflict)", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("DeleteGlobalCollege", mock.Anything, collegeID).
			Return(errors.New("cannot delete this college: one or more personal college applications still reference it"))

		handler := NewHandler(mockRepo)
		err := handler.DeleteGlobalCollege(ctx, collegeID)

		assert.Error(t, err)
	})
}
