package essays

import (
	"context"
	"errors"
	"net/http"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Unit tests for the UpdateStatus handler logic without HTTP or database dependencies.
func TestHandler_UpdateStatus(t *testing.T) {
	t.Parallel()

	essayID := uuid.New()
	studentID := uuid.New()

	input := &models.UpdateStatusInput{
		EssayID: essayID,
		Body: models.UpdateStatusBody{
			Status: models.Submitted,
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		updatedEssay := &models.Essays{
			ID:     essayID,
			Status: models.Submitted,
		}

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssayStudentID", mock.Anything, essayID).Return(studentID, nil)
		mockRepo.On("UpdateStatus", mock.Anything, essayID, models.Submitted).Return(updatedEssay, nil)

		handler := NewHandler(mockRepo, realAuthorizer(), nil)
		res, err := handler.UpdateStatus(studentRoleCtx(studentID), input)

		assert.NoError(t, err)
		assert.Equal(t, &models.UpdateStatusOutput{
			Body: models.EssayBody{Essay: updatedEssay},
		}, res)
	})

	t.Run("student cannot update another student's essay", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssayStudentID", mock.Anything, essayID).Return(uuid.New(), nil)

		handler := NewHandler(mockRepo, realAuthorizer(), nil)
		res, err := handler.UpdateStatus(studentRoleCtx(studentID), input)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.GetStatus())
		mockRepo.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("essay lookup error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssayStudentID", mock.Anything, essayID).
			Return(uuid.Nil, errs.NotFound("essay", "id", essayID.String()))

		handler := NewHandler(mockRepo, nil, nil)
		res, err := handler.UpdateStatus(context.Background(), input)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.GetStatus())
		mockRepo.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssayStudentID", mock.Anything, essayID).Return(studentID, nil)
		mockRepo.On("UpdateStatus", mock.Anything, essayID, models.Submitted).
			Return((*models.Essays)(nil), errors.New("database error"))

		handler := NewHandler(mockRepo, nil, nil)
		res, err := handler.UpdateStatus(context.Background(), input)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
