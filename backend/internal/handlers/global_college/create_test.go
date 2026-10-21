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

func TestHandler_CreateGlobalCollege(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Northeastern Boston Campus",
		SchoolLocation: "Boston MA",
	}

	t.Run("no near-duplicate: creates the college", func(t *testing.T) {
		t.Parallel()

		expected := &models.GlobalCollege{ID: 1, SchoolName: input.SchoolName, SchoolLocation: input.SchoolLocation}

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("FindNearDuplicateGlobalColleges", mock.Anything, input.SchoolName, input.SchoolLocation, mock.Anything).
			Return([]models.GlobalCollege{}, nil)
		mockRepo.On("CreateGlobalCollege", mock.Anything, input).Return(expected, nil)

		res, err := NewHandler(mockRepo).CreateGlobalCollege(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("near-duplicate found: rejected with 409, nothing created", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("FindNearDuplicateGlobalColleges", mock.Anything, input.SchoolName, input.SchoolLocation, mock.Anything).
			Return([]models.GlobalCollege{{ID: 7, SchoolName: "Northeastern University - Boston", SchoolLocation: "Boston, MA"}}, nil)
		// CreateGlobalCollege is not mocked: calling it would fail the test.

		res, err := NewHandler(mockRepo).CreateGlobalCollege(ctx, input)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusConflict, httpErr.Code)
	})

	t.Run("confirm_duplicate skips the near-duplicate check", func(t *testing.T) {
		t.Parallel()

		confirmed := input
		confirmed.ConfirmDuplicate = true
		expected := &models.GlobalCollege{ID: 2, SchoolName: input.SchoolName, SchoolLocation: input.SchoolLocation}

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		// FindNearDuplicateGlobalColleges is not mocked: calling it would fail the test.
		mockRepo.On("CreateGlobalCollege", mock.Anything, confirmed).Return(expected, nil)

		res, err := NewHandler(mockRepo).CreateGlobalCollege(ctx, confirmed)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("near-duplicate lookup error propagates, nothing created", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("FindNearDuplicateGlobalColleges", mock.Anything, input.SchoolName, input.SchoolLocation, mock.Anything).
			Return(nil, errors.New("database error"))

		res, err := NewHandler(mockRepo).CreateGlobalCollege(ctx, input)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})

	t.Run("create error propagates (e.g. exact duplicate)", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("FindNearDuplicateGlobalColleges", mock.Anything, input.SchoolName, input.SchoolLocation, mock.Anything).
			Return([]models.GlobalCollege{}, nil)
		mockRepo.On("CreateGlobalCollege", mock.Anything, input).
			Return(nil, errs.Conflict("global college", "school_name/school_location", input.SchoolName+" / "+input.SchoolLocation))

		res, err := NewHandler(mockRepo).CreateGlobalCollege(ctx, input)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusConflict, httpErr.Code)
	})
}
