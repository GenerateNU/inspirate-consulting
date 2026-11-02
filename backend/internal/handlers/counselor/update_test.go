package counselor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func TestHandler_UpdateCounselor(t *testing.T) {
	t.Parallel()

	input := models.UpdateCounselorBody{}

	t.Run("updates the logged-in counselor", func(t *testing.T) {
		t.Parallel()

		ctx := auth.WithCounselorID(context.Background(), testCounselorID)

		repo := mocks.NewCounselorRepository(t)
		repo.On("UpdateCounselor", ctx, testCounselorID, input).Return(&models.Counselor{ID: testCounselorID}, nil)

		res, err := NewHandler(repo).UpdateCounselor(ctx, input)

		require.NoError(t, err)
		assert.Equal(t, testCounselorID, res.ID)
	})

	t.Run("caller without a counselor profile is unauthorized", func(t *testing.T) {
		t.Parallel()

		repo := mocks.NewCounselorRepository(t) // UpdateCounselor should not not be called

		res, err := NewHandler(repo).UpdateCounselor(context.Background(), input)

		assert.Nil(t, res)
		assert.Equal(t, errs.Unauthorized(), err)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		ctx := auth.WithCounselorID(context.Background(), testCounselorID)
		notFound := errs.NotFound("counselor", "id", testCounselorID.String())

		repo := mocks.NewCounselorRepository(t)
		repo.On("UpdateCounselor", ctx, testCounselorID, input).Return(nil, notFound)

		res, err := NewHandler(repo).UpdateCounselor(ctx, input)

		assert.Nil(t, res)
		assert.Equal(t, notFound, err)
	})
}
