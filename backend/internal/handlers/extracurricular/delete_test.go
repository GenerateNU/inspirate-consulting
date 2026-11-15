package extracurricular

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

// Unit tests for the DeleteExtracurricular handler.
func TestHandler_DeleteExtracurricular(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	const id int64 = 21
	existing := &models.ExtracurricularDetail{
		Extracurricular: models.Extracurricular{ID: id, StudentID: ownerID.String(), Name: "Debate"},
	}

	t.Run("owner deletes using the row's student id", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(existing, nil)
		mockRepo.On("DeleteExtracurricular", mock.Anything, id, ownerID.String()).Return(nil)

		err := NewHandler(mockRepo, realAuthorizer()).DeleteExtracurricular(studentRoleCtx(ownerID), id)

		assert.NoError(t, err)
	})

	t.Run("nil authorizer (test mode) uses the owner from the row, not the caller", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(existing, nil)
		mockRepo.On("DeleteExtracurricular", mock.Anything, id, ownerID.String()).Return(nil)

		err := NewHandler(mockRepo, nil).DeleteExtracurricular(context.Background(), id)

		assert.NoError(t, err)
	})

	t.Run("another student is forbidden and nothing is deleted", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(existing, nil)

		err := NewHandler(mockRepo, realAuthorizer()).DeleteExtracurricular(studentRoleCtx(uuid.New()), id)

		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		mockRepo.AssertNotCalled(t, "DeleteExtracurricular", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("not found on lookup skips the delete", func(t *testing.T) {
		t.Parallel()

		notFound := errs.NotFound("extracurricular", "id", id)
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(nil, notFound)

		err := NewHandler(mockRepo, realAuthorizer()).DeleteExtracurricular(studentRoleCtx(ownerID), id)

		assert.Equal(t, notFound, err)
		mockRepo.AssertNotCalled(t, "DeleteExtracurricular", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("not found on delete (concurrent removal) is propagated", func(t *testing.T) {
		t.Parallel()

		notFound := errs.NotFound("extracurricular", "id", id)
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(existing, nil)
		mockRepo.On("DeleteExtracurricular", mock.Anything, id, ownerID.String()).Return(notFound)

		err := NewHandler(mockRepo, realAuthorizer()).DeleteExtracurricular(studentRoleCtx(ownerID), id)

		assert.Equal(t, notFound, err)
	})

	t.Run("lookup repository error is propagated", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(nil, errors.New("database error"))

		err := NewHandler(mockRepo, realAuthorizer()).DeleteExtracurricular(studentRoleCtx(ownerID), id)

		assert.EqualError(t, err, "database error")
		mockRepo.AssertNotCalled(t, "DeleteExtracurricular", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("delete repository error is propagated", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(existing, nil)
		mockRepo.On("DeleteExtracurricular", mock.Anything, id, ownerID.String()).Return(errors.New("database error"))

		err := NewHandler(mockRepo, realAuthorizer()).DeleteExtracurricular(studentRoleCtx(ownerID), id)

		assert.EqualError(t, err, "database error")
	})
}
