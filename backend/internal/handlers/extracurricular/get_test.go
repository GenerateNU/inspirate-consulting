package extracurricular

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Unit tests for the GetExtracurricular handler.
func TestHandler_GetExtracurricular(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	const id int64 = 11

	detail := &models.ExtracurricularDetail{
		Extracurricular: models.Extracurricular{
			ID:                  id,
			StudentID:           ownerID.String(),
			Name:                "Debate",
			HoursPerWeek:        ptr(6.0),
			ReadyForApplication: true,
		},
		Achievements: []models.ExtracurricularAchievement{
			{ID: 1, ExtracurricularID: id, EntryDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC), Body: "Won state"},
		},
		Notes: []models.ExtracurricularNote{
			{ID: 2, ExtracurricularID: id, Body: "Ask for rec letter"},
		},
	}

	t.Run("owner gets the detail with achievements and notes", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(detail, nil)

		res, err := NewHandler(mockRepo, realAuthorizer()).GetExtracurricular(studentRoleCtx(ownerID), id)

		require.NoError(t, err)
		assert.Equal(t, detail, res)
		assert.Len(t, res.Achievements, 1)
		assert.Len(t, res.Notes, 1)
	})

	t.Run("nil authorizer (test mode) allows access", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(detail, nil)

		res, err := NewHandler(mockRepo, nil).GetExtracurricular(context.Background(), id)

		require.NoError(t, err)
		assert.Equal(t, detail, res)
	})

	t.Run("another student is forbidden", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(detail, nil)

		res, err := NewHandler(mockRepo, realAuthorizer()).GetExtracurricular(studentRoleCtx(uuid.New()), id)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("caller without a role is forbidden", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(detail, nil)

		res, err := NewHandler(mockRepo, realAuthorizer()).GetExtracurricular(context.Background(), id)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("not found is propagated", func(t *testing.T) {
		t.Parallel()

		notFound := errs.NotFound("extracurricular", "id", id)
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(nil, notFound)

		res, err := NewHandler(mockRepo, realAuthorizer()).GetExtracurricular(studentRoleCtx(ownerID), id)

		assert.Nil(t, res)
		assert.Equal(t, notFound, err)
	})

	t.Run("corrupt owner id is a bad request, not a leak", func(t *testing.T) {
		t.Parallel()

		bad := &models.ExtracurricularDetail{Extracurricular: models.Extracurricular{ID: id, StudentID: "not-a-uuid"}}
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(bad, nil)

		res, err := NewHandler(mockRepo, realAuthorizer()).GetExtracurricular(studentRoleCtx(ownerID), id)

		assert.Nil(t, res)
		assert.Error(t, err)
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, id).Return(nil, errors.New("database error"))

		res, err := NewHandler(mockRepo, realAuthorizer()).GetExtracurricular(studentRoleCtx(ownerID), id)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
