package globalcollege

import (
	"context"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_DeleteGlobalCollege(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("DeleteGlobalCollege", mock.Anything, int64(1)).Return(nil)

		err := NewHandler(mockRepo).DeleteGlobalCollege(ctx, 1)

		assert.NoError(t, err)
	})

	t.Run("not found propagates", func(t *testing.T) {
		t.Parallel()

		notFound := errs.NotFound("global college", "id", int64(999))

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("DeleteGlobalCollege", mock.Anything, int64(999)).Return(notFound)

		err := NewHandler(mockRepo).DeleteGlobalCollege(ctx, 999)

		assert.Equal(t, notFound, err)
	})

	t.Run("still referenced propagates", func(t *testing.T) {
		t.Parallel()

		inUse := errs.BadRequest("cannot delete: this college is still in use")

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("DeleteGlobalCollege", mock.Anything, int64(1)).Return(inUse)

		err := NewHandler(mockRepo).DeleteGlobalCollege(ctx, 1)

		assert.Equal(t, inUse, err)
	})
}