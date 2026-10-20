package student

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// expectTx makes the mock run the callback the handler passes to WithTx, so
// the orchestration inside it is exercised. The mocked queries ignore db.
func expectTx(mockRepo *mocks.EssayReviewRepository) {
	mockRepo.On("WithTx", mock.Anything, mock.Anything).
		Return(func(ctx context.Context, fn func(db dbinterface.QueryInterface) error) error {
			return fn(nil)
		})
}

// Unit tests for the SetReviewBalance handler logic without HTTP or database dependencies.
func TestHandler_SetReviewBalance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("negative balance is rejected before any database work", func(t *testing.T) {
		t.Parallel()

		mockEssayRepo := mocks.NewEssayReviewRepository(t)
		mockStudentRepo := mocks.NewStudentRepository(t)
		handler := NewHandler(mockEssayRepo, mockStudentRepo, nil)

		res, err := handler.SetReviewBalance(ctx, uuid.New(), models.SetStudentReviewBalanceRequestBody{
			ReviewBalance: -1,
		})

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 400, err.(errs.HTTPError).Code)
		mockEssayRepo.AssertNotCalled(t, "WithTx", mock.Anything, mock.Anything)
	})

	t.Run("raising the balance records a positive adjustment", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()

		mockEssayRepo := mocks.NewEssayReviewRepository(t)
		mockStudentRepo := mocks.NewStudentRepository(t)
		expectTx(mockEssayRepo)
		mockEssayRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(&models.Student{ID: studentID, ReviewBalance: 2}, nil)
		mockEssayRepo.On("SetStudentBalance", mock.Anything, mock.Anything, studentID, 7).Return(nil)
		mockEssayRepo.On("InsertAdjustment", mock.Anything, mock.Anything, studentID, 5).Return(nil)

		handler := NewHandler(mockEssayRepo, mockStudentRepo, nil)
		res, err := handler.SetReviewBalance(ctx, studentID, models.SetStudentReviewBalanceRequestBody{
			ReviewBalance: 7,
		})

		require.NoError(t, err)
		assert.Equal(t, 7, res.ReviewBalance)
	})

	t.Run("lowering the balance records a negative adjustment", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()

		mockEssayRepo := mocks.NewEssayReviewRepository(t)
		mockStudentRepo := mocks.NewStudentRepository(t)
		expectTx(mockEssayRepo)
		mockEssayRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(&models.Student{ID: studentID, ReviewBalance: 5}, nil)
		mockEssayRepo.On("SetStudentBalance", mock.Anything, mock.Anything, studentID, 1).Return(nil)
		mockEssayRepo.On("InsertAdjustment", mock.Anything, mock.Anything, studentID, -4).Return(nil)

		handler := NewHandler(mockEssayRepo, mockStudentRepo, nil)
		res, err := handler.SetReviewBalance(ctx, studentID, models.SetStudentReviewBalanceRequestBody{
			ReviewBalance: 1,
		})

		require.NoError(t, err)
		assert.Equal(t, 1, res.ReviewBalance)
	})

	t.Run("setting the same value writes nothing", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()

		mockEssayRepo := mocks.NewEssayReviewRepository(t)
		mockStudentRepo := mocks.NewStudentRepository(t)
		expectTx(mockEssayRepo)
		mockEssayRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(&models.Student{ID: studentID, ReviewBalance: 4}, nil)

		handler := NewHandler(mockEssayRepo, mockStudentRepo, nil)
		res, err := handler.SetReviewBalance(ctx, studentID, models.SetStudentReviewBalanceRequestBody{
			ReviewBalance: 4,
		})

		require.NoError(t, err)
		assert.Equal(t, 4, res.ReviewBalance)
		// a zero delta would be a meaningless ledger row
		mockEssayRepo.AssertNotCalled(t, "SetStudentBalance", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		mockEssayRepo.AssertNotCalled(t, "InsertAdjustment", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("unknown student propagates the repository error", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()

		mockEssayRepo := mocks.NewEssayReviewRepository(t)
		mockStudentRepo := mocks.NewStudentRepository(t)
		expectTx(mockEssayRepo)
		mockEssayRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(nil, errs.NotFound("student", "id", studentID.String()))

		handler := NewHandler(mockEssayRepo, mockStudentRepo, nil)
		res, err := handler.SetReviewBalance(ctx, studentID, models.SetStudentReviewBalanceRequestBody{
			ReviewBalance: 3,
		})

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 404, err.(errs.HTTPError).Code)
	})
}
