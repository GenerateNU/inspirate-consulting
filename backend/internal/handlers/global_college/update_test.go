package globalcollege

import (
	"context"
	"errors"
	"net/http"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_UpdateGlobalCollege(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := int64(1)
	input := models.UpdateGlobalCollegeRequestBody{
		SchoolName:     "Northeastern Boston Campus",
		SchoolLocation: "Boston MA",
	}

	t.Run("no near-duplicate: updates the college", func(t *testing.T) {
		t.Parallel()

		expected := &models.GlobalCollege{ID: id, SchoolName: input.SchoolName, SchoolLocation: input.SchoolLocation}

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("FindNearDuplicateGlobalColleges", mock.Anything, input.SchoolName, input.SchoolLocation, nearDuplicateThreshold).
			Return([]models.GlobalCollege{}, nil)
		mockRepo.On("UpdateGlobalCollege", mock.Anything, id, input).Return(expected, nil)

		res, err := NewHandler(mockRepo).UpdateGlobalCollege(ctx, id, input)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("only match is the college being updated: updates (no self false-positive)", func(t *testing.T) {
		t.Parallel()

		expected := &models.GlobalCollege{ID: id, SchoolName: input.SchoolName, SchoolLocation: input.SchoolLocation}

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("FindNearDuplicateGlobalColleges", mock.Anything, input.SchoolName, input.SchoolLocation, nearDuplicateThreshold).
			Return([]models.GlobalCollege{{ID: id, SchoolName: input.SchoolName}}, nil)
		mockRepo.On("UpdateGlobalCollege", mock.Anything, id, input).Return(expected, nil)

		res, err := NewHandler(mockRepo).UpdateGlobalCollege(ctx, id, input)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("near-duplicate of another college: rejected with 409, nothing updated", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("FindNearDuplicateGlobalColleges", mock.Anything, input.SchoolName, input.SchoolLocation, nearDuplicateThreshold).
			Return([]models.GlobalCollege{
				{ID: id, SchoolName: input.SchoolName},
				{ID: 7, SchoolName: "Northeastern University - Boston"},
			}, nil)
		// UpdateGlobalCollege is not mocked: calling it would fail the test.

		res, err := NewHandler(mockRepo).UpdateGlobalCollege(ctx, id, input)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusConflict, httpErr.Code)
	})

	t.Run("confirm_duplicate skips the near-duplicate check", func(t *testing.T) {
		t.Parallel()

		confirmed := input
		confirmed.ConfirmDuplicate = true
		expected := &models.GlobalCollege{ID: id, SchoolName: input.SchoolName, SchoolLocation: input.SchoolLocation}

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		// FindNearDuplicateGlobalColleges is not mocked: calling it would fail the test.
		mockRepo.On("UpdateGlobalCollege", mock.Anything, id, confirmed).Return(expected, nil)

		res, err := NewHandler(mockRepo).UpdateGlobalCollege(ctx, id, confirmed)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("near-duplicate lookup error propagates, nothing updated", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("FindNearDuplicateGlobalColleges", mock.Anything, input.SchoolName, input.SchoolLocation, nearDuplicateThreshold).
			Return(nil, errors.New("database error"))

		res, err := NewHandler(mockRepo).UpdateGlobalCollege(ctx, id, input)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})

	t.Run("update error propagates (e.g. not found)", func(t *testing.T) {
		t.Parallel()

		notFound := errs.NotFound("global college", "id", id)

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("FindNearDuplicateGlobalColleges", mock.Anything, input.SchoolName, input.SchoolLocation, nearDuplicateThreshold).
			Return([]models.GlobalCollege{}, nil)
		mockRepo.On("UpdateGlobalCollege", mock.Anything, id, input).Return(nil, notFound)

		res, err := NewHandler(mockRepo).UpdateGlobalCollege(ctx, id, input)

		assert.Nil(t, res)
		assert.Equal(t, notFound, err)
	})
}
