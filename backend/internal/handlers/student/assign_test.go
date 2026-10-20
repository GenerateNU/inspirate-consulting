package student

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func assignInput(userID, counselorID uuid.UUID) *models.AssignCounselorInput {
	input := &models.AssignCounselorInput{}
	input.Body.UserID = userID
	input.Body.CounselorID = counselorID
	return input
}

// Unit tests for the AssignCounselor handler logic without HTTP or database dependencies.
func TestHandler_AssignCounselor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("assigns the counselor and returns the updated student", func(t *testing.T) {
		t.Parallel()

		studentUserID := uuid.New()
		counselorID := uuid.New()
		updated := &models.Student{ID: uuid.New(), UserID: studentUserID, Year: models.Junior, CounselorID: &counselorID}

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("AssignCounselor", ctx, studentUserID, counselorID).Return(updated, nil)

		handler := NewHandler(nil, mockRepo, nil)
		res, err := handler.AssignCounselor(ctx, assignInput(studentUserID, counselorID))

		require.NoError(t, err)
		assert.Equal(t, updated, res)
		require.NotNil(t, res.CounselorID)
		assert.Equal(t, counselorID, *res.CounselorID)
	})

	t.Run("unknown student returns not found", func(t *testing.T) {
		t.Parallel()

		studentUserID := uuid.New()
		counselorID := uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("AssignCounselor", ctx, studentUserID, counselorID).
			Return(nil, errs.NotFound("student", "user_id", studentUserID))

		handler := NewHandler(nil, mockRepo, nil)
		res, err := handler.AssignCounselor(ctx, assignInput(studentUserID, counselorID))

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 404, err.(errs.HTTPError).Code)
	})

	t.Run("unknown counselor returns not found", func(t *testing.T) {
		t.Parallel()

		studentUserID := uuid.New()
		counselorID := uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("AssignCounselor", ctx, studentUserID, counselorID).
			Return(nil, errs.NotFound("counselor", "id", counselorID))

		handler := NewHandler(nil, mockRepo, nil)
		res, err := handler.AssignCounselor(ctx, assignInput(studentUserID, counselorID))

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 404, err.(errs.HTTPError).Code)
	})

	t.Run("assigning a counselor to another counselor fails", func(t *testing.T) {
		t.Parallel()

		// A counselor's user id has no student row, so the repository finds no student to update.
		counselorUserID := uuid.New()
		counselorID := uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("AssignCounselor", ctx, counselorUserID, counselorID).
			Return(nil, errs.NotFound("student", "user_id", counselorUserID))

		handler := NewHandler(nil, mockRepo, nil)
		res, err := handler.AssignCounselor(ctx, assignInput(counselorUserID, counselorID))

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 404, err.(errs.HTTPError).Code)
	})
}
