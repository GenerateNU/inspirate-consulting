package personalcollegeapplication

import (
	"context"
	"errors"
	"testing"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func intPtr(i int) *int { return &i }

func TestHandler_UpdatePersonalCollegeApplicationRank(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	ctx := auth.WithStudentID(context.Background(), studentID)
	applicationID := int64(42)

	t.Run("success : set a rank", func(t *testing.T) {
		t.Parallel()

		expected := []models.PersonalCollegeApplication{
			{ID: applicationID, StudentID: studentID, Rank: intPtr(1)},
			{ID: 43, StudentID: studentID, Rank: intPtr(2)},
		}

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("UpdateApplicationRank", mock.Anything, studentID, applicationID, intPtr(1)).Return(expected, nil)

		handler := NewHandler(mockPCARepo, nil)
		res, err := handler.UpdatePersonalCollegeApplicationRank(ctx, applicationID, models.UpdatePersonalCollegeApplicationRankRequestBody{
			Rank: intPtr(1),
		})

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("success : unrank (nil)", func(t *testing.T) {
		t.Parallel()

		expected := []models.PersonalCollegeApplication{
			{ID: applicationID, StudentID: studentID, Rank: nil},
		}

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("UpdateApplicationRank", mock.Anything, studentID, applicationID, (*int)(nil)).Return(expected, nil)

		handler := NewHandler(mockPCARepo, nil)
		res, err := handler.UpdatePersonalCollegeApplicationRank(ctx, applicationID, models.UpdatePersonalCollegeApplicationRankRequestBody{
			Rank: nil,
		})

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("repository error propagates (e.g. not found)", func(t *testing.T) {
		t.Parallel()

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("UpdateApplicationRank", mock.Anything, studentID, applicationID, intPtr(1)).
			Return(nil, errors.New("personal college application not found"))

		handler := NewHandler(mockPCARepo, nil)
		res, err := handler.UpdatePersonalCollegeApplicationRank(ctx, applicationID, models.UpdatePersonalCollegeApplicationRankRequestBody{
			Rank: intPtr(1),
		})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("missing student id returns error", func(t *testing.T) {
		t.Parallel()

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		// UpdateApplicationRank is not called without a student id

		handler := NewHandler(mockPCARepo, nil)
		res, err := handler.UpdatePersonalCollegeApplicationRank(context.Background(), applicationID, models.UpdatePersonalCollegeApplicationRankRequestBody{
			Rank: intPtr(1),
		})

		assert.ErrorIs(t, err, auth.ErrNoStudentID)
		assert.Nil(t, res)
	})
}
