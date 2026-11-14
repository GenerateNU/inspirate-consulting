package activity

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
)

var testCounselorID = uuid.MustParse("00000000-0000-0000-0000-000000000003")

// Unit tests for the GetActivities handler logic.
// The counselor ID is set on the context the same way the auth middleware does.
func TestHandler_GetActivities(t *testing.T) {
	t.Parallel()

	ctx := auth.WithCounselorID(context.Background(), testCounselorID)

	t.Run("missing counselor id returns error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewActivityRepository(t)
		// no repository methods are called without a counselor id

		handler := NewHandler(mockRepo)
		res, err := handler.GetActivities(context.Background(), &models.GetActivitiesInput{})

		assert.ErrorIs(t, err, auth.ErrNoCounselorID)
		assert.Nil(t, res)
	})

	t.Run("passes the caller's counselor id and no student filter through to the repository", func(t *testing.T) {
		t.Parallel()

		want := []models.Activity{{ID: uuid.New(), Action: models.ActionAddedSchool}}

		mockRepo := mocks.NewActivityRepository(t)
		mockRepo.On("GetActivitiesForCounselor", ctx, testCounselorID, (*uuid.UUID)(nil)).Return(want, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.GetActivities(ctx, &models.GetActivitiesInput{})

		assert.NoError(t, err)
		assert.Equal(t, want, res)
	})

	t.Run("passes a student_id filter through to the repository", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()
		want := []models.Activity{{ID: uuid.New(), StudentID: studentID, Action: models.ActionTaskPastDue}}

		mockRepo := mocks.NewActivityRepository(t)
		mockRepo.On("GetActivitiesForCounselor", ctx, testCounselorID, &studentID).Return(want, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.GetActivities(ctx, &models.GetActivitiesInput{StudentID: studentID})

		assert.NoError(t, err)
		assert.Equal(t, want, res)
	})
}
