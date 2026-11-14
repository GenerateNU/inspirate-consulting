package extracurricular

import (
	"context"
	"errors"
	"testing"
	"time"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Unit tests for the CreateExtracurricular handler without HTTP or database dependencies.
func TestHandler_CreateExtracurricular(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	userID := uuid.New()

	input := &models.CreateExtracurricularInput{
		Body: models.CreateExtracurricularRequest{
			Name:                   "Varsity Debate",
			Status:                 models.ExtracurricularStatusDoing,
			Type:                   models.ExtracurricularTypeInvestment,
			Description:            ptr("Regional tournaments"),
			StartDate:              "2024-09-01",
			Organization:           "Lincoln High",
			HoursPerWeek:           ptr(7.5),
			WeeksPerYear:           ptr(int16(36)),
			Goals:                  ptr("Qualify for nationals"),
			ReadyForApplication:    ptr(true),
			ApplicationDescription: ptr("Led the team to a state title."),
		},
	}
	expected := &models.CreateExtracurricularOutput{
		Body: models.Extracurricular{
			ID:                     1,
			StudentID:              studentID.String(),
			UserID:                 userID.String(),
			Name:                   "Varsity Debate",
			Status:                 models.ExtracurricularStatusDoing,
			Type:                   models.ExtracurricularTypeInvestment,
			Description:            ptr("Regional tournaments"),
			StartDate:              time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
			Organization:           "Lincoln High",
			HoursPerWeek:           ptr(7.5),
			WeeksPerYear:           ptr(int16(36)),
			Goals:                  ptr("Qualify for nationals"),
			ReadyForApplication:    true,
			ApplicationDescription: ptr("Led the team to a state title."),
		},
	}

	t.Run("success passes identity from context and the full body", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("CreateExtracurricular", mock.Anything, studentID.String(), userID.String(), *input).Return(expected, nil)

		res, err := NewHandler(mockRepo, nil).CreateExtracurricular(studentCtx(studentID, userID), input)

		require.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("new fields reach the repository untouched", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("CreateExtracurricular", mock.Anything, studentID.String(), userID.String(),
			mock.MatchedBy(func(in models.CreateExtracurricularInput) bool {
				b := in.Body
				return *b.HoursPerWeek == 7.5 &&
					*b.WeeksPerYear == 36 &&
					*b.Goals == "Qualify for nationals" &&
					*b.ReadyForApplication &&
					*b.ApplicationDescription == "Led the team to a state title."
			})).Return(expected, nil)

		_, err := NewHandler(mockRepo, nil).CreateExtracurricular(studentCtx(studentID, userID), input)
		require.NoError(t, err)
	})

	t.Run("minimal body with no optional fields", func(t *testing.T) {
		t.Parallel()

		minimal := &models.CreateExtracurricularInput{
			Body: models.CreateExtracurricularRequest{
				Name:         "Chess",
				Status:       models.ExtracurricularStatusHaveDone,
				Type:         models.ExtracurricularTypeMaintenance,
				StartDate:    "2023-01-01",
				Organization: "Chess Club",
			},
		}
		out := &models.CreateExtracurricularOutput{Body: models.Extracurricular{ID: 2, Name: "Chess"}}

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("CreateExtracurricular", mock.Anything, studentID.String(), userID.String(), *minimal).Return(out, nil)

		res, err := NewHandler(mockRepo, nil).CreateExtracurricular(studentCtx(studentID, userID), minimal)

		require.NoError(t, err)
		assert.Nil(t, res.Body.Description)
		assert.Nil(t, res.Body.HoursPerWeek)
		assert.False(t, res.Body.ReadyForApplication)
	})

	t.Run("missing student id is unauthorized", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		ctx := auth.WithUserID(context.Background(), userID)

		res, err := NewHandler(mockRepo, nil).CreateExtracurricular(ctx, input)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, 401, httpErr.Code)
		mockRepo.AssertNotCalled(t, "CreateExtracurricular")
	})

	t.Run("missing user id is unauthorized", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		ctx := auth.WithStudentID(context.Background(), studentID)

		res, err := NewHandler(mockRepo, nil).CreateExtracurricular(ctx, input)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, 401, httpErr.Code)
		mockRepo.AssertNotCalled(t, "CreateExtracurricular")
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("CreateExtracurricular", mock.Anything, studentID.String(), userID.String(), *input).
			Return(nil, errors.New("database error"))

		res, err := NewHandler(mockRepo, nil).CreateExtracurricular(studentCtx(studentID, userID), input)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
