package essayreview

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// Unit tests for the RequestEssayReview handler logic without HTTP or database dependencies.
func TestHandler_RequestEssayReview(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("amount below one is rejected before any database work", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayReviewRepository(t)
		handler := NewHandler(mockRepo)

		res, err := handler.RequestEssayReview(ctx, models.RequestEssayReviewRequestBody{
			StudentID: uuid.New(), EssayID: uuid.New(), Amount: 0,
		})

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 400, err.(errs.HTTPError).Code)
		mockRepo.AssertNotCalled(t, "WithTx", mock.Anything, mock.Anything)
	})

	t.Run("essay with an open review is rejected and no balance moves", func(t *testing.T) {
		t.Parallel()

		studentID, essayID, openID := uuid.New(), uuid.New(), uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(&models.Student{ID: studentID, ReviewBalance: 10}, nil)
		mockRepo.On("FindOpenReviewForEssay", mock.Anything, mock.Anything, essayID).
			Return(&openID, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.RequestEssayReview(ctx, models.RequestEssayReviewRequestBody{
			StudentID: studentID, EssayID: essayID, Amount: 1,
		})

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 409, err.(errs.HTTPError).Code)
		assert.Contains(t, err.(errs.HTTPError).Message, openID.String())
		mockRepo.AssertNotCalled(t, "AdjustStudentBalance", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		mockRepo.AssertNotCalled(t, "InsertSpend", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("insufficient balance is rejected and no balance moves", func(t *testing.T) {
		t.Parallel()

		studentID, essayID := uuid.New(), uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(&models.Student{ID: studentID, ReviewBalance: 2}, nil)
		mockRepo.On("FindOpenReviewForEssay", mock.Anything, mock.Anything, essayID).
			Return(nil, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.RequestEssayReview(ctx, models.RequestEssayReviewRequestBody{
			StudentID: studentID, EssayID: essayID, Amount: 5,
		})

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 400, err.(errs.HTTPError).Code)
		mockRepo.AssertNotCalled(t, "AdjustStudentBalance", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		mockRepo.AssertNotCalled(t, "InsertSpend", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("success debits the balance then records the spend", func(t *testing.T) {
		t.Parallel()

		studentID, essayID := uuid.New(), uuid.New()
		expectedEntity := transactionRow(func(row *models.EssayReviewTransaction) {
			row.StudentID = studentID
			row.EssayID = &essayID
		})

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(&models.Student{ID: studentID, ReviewBalance: 10}, nil)
		mockRepo.On("FindOpenReviewForEssay", mock.Anything, mock.Anything, essayID).
			Return(nil, nil)
		mockRepo.On("AdjustStudentBalance", mock.Anything, mock.Anything, studentID, -3).Return(nil)
		mockRepo.On("InsertSpend", mock.Anything, mock.Anything, studentID, essayID, 3).
			Return(expectedEntity, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.RequestEssayReview(ctx, models.RequestEssayReviewRequestBody{
			StudentID: studentID, EssayID: essayID, Amount: 3,
		})

		assert.NoError(t, err)
		assert.Equal(t, expectedEntity, res)
	})

	t.Run("unknown student propagates the repository error", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(nil, errs.NotFound("student", "id", studentID.String()))

		handler := NewHandler(mockRepo)
		res, err := handler.RequestEssayReview(ctx, models.RequestEssayReviewRequestBody{
			StudentID: studentID, EssayID: uuid.New(), Amount: 1,
		})

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 404, err.(errs.HTTPError).Code)
	})
}
