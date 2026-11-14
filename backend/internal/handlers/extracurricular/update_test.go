package extracurricular

import (
	"errors"
	"net/http"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Unit tests for the UpdateExtracurricular handler.
func TestHandler_UpdateExtracurricular(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	userID := uuid.New()
	const id int64 = 42

	t.Run("success updating only new fields", func(t *testing.T) {
		t.Parallel()

		input := &models.UpdateExtracurricularInput{
			ID: id,
			Body: models.UpdateExtracurricularRequest{
				HoursPerWeek:           ptr(12.5),
				WeeksPerYear:           ptr(int16(48)),
				Goals:                  ptr("Become captain"),
				ReadyForApplication:    ptr(true),
				ApplicationDescription: ptr("Updated blurb"),
			},
		}
		expected := &models.Extracurricular{
			ID:                     id,
			StudentID:              studentID.String(),
			Name:                   "Debate",
			HoursPerWeek:           ptr(12.5),
			WeeksPerYear:           ptr(int16(48)),
			Goals:                  ptr("Become captain"),
			ReadyForApplication:    true,
			ApplicationDescription: ptr("Updated blurb"),
		}

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricularStudentID", mock.Anything, id).Return(studentID, nil)
		mockRepo.On("UpdateExtracurricular", mock.Anything, id, *input).Return(expected, nil)

		res, err := NewHandler(mockRepo, nil).UpdateExtracurricular(studentCtx(studentID, userID), id, input)

		require.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("explicit false for ready_for_application is forwarded, not dropped", func(t *testing.T) {
		t.Parallel()

		input := &models.UpdateExtracurricularInput{
			ID:   id,
			Body: models.UpdateExtracurricularRequest{ReadyForApplication: ptr(false)},
		}

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricularStudentID", mock.Anything, id).Return(studentID, nil)
		mockRepo.On("UpdateExtracurricular", mock.Anything, id, mock.MatchedBy(func(in models.UpdateExtracurricularInput) bool {
			return in.Body.ReadyForApplication != nil && !*in.Body.ReadyForApplication
		})).Return(&models.Extracurricular{ID: id}, nil)

		_, err := NewHandler(mockRepo, nil).UpdateExtracurricular(studentCtx(studentID, userID), id, input)
		require.NoError(t, err)
	})

	t.Run("empty body is forwarded", func(t *testing.T) {
		t.Parallel()

		input := &models.UpdateExtracurricularInput{ID: id}
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricularStudentID", mock.Anything, id).Return(studentID, nil)
		mockRepo.On("UpdateExtracurricular", mock.Anything, id, *input).Return(&models.Extracurricular{ID: id}, nil)

		res, err := NewHandler(mockRepo, nil).UpdateExtracurricular(studentCtx(studentID, userID), id, input)

		require.NoError(t, err)
		assert.Equal(t, id, res.ID)
	})

	t.Run("owner passes the real authorizer", func(t *testing.T) {
		t.Parallel()

		input := &models.UpdateExtracurricularInput{ID: id, Body: models.UpdateExtracurricularRequest{Goals: ptr("x")}}
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricularStudentID", mock.Anything, id).Return(studentID, nil)
		mockRepo.On("UpdateExtracurricular", mock.Anything, id, *input).Return(&models.Extracurricular{ID: id}, nil)

		res, err := NewHandler(mockRepo, realAuthorizer()).UpdateExtracurricular(studentRoleCtx(studentID), id, input)

		require.NoError(t, err)
		assert.Equal(t, id, res.ID)
	})

	t.Run("another student is forbidden", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricularStudentID", mock.Anything, id).Return(studentID, nil)
		input := &models.UpdateExtracurricularInput{ID: id, Body: models.UpdateExtracurricularRequest{Goals: ptr("x")}}

		res, err := NewHandler(mockRepo, realAuthorizer()).UpdateExtracurricular(studentRoleCtx(uuid.New()), id, input)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
		mockRepo.AssertNotCalled(t, "UpdateExtracurricular")
	})

	t.Run("owner lookup not found is propagated", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricularStudentID", mock.Anything, id).Return(uuid.Nil, pgx.ErrNoRows)
		input := &models.UpdateExtracurricularInput{ID: id, Body: models.UpdateExtracurricularRequest{Goals: ptr("x")}}

		res, err := NewHandler(mockRepo, nil).UpdateExtracurricular(studentCtx(studentID, userID), id, input)

		assert.Nil(t, res)
		assert.ErrorIs(t, err, pgx.ErrNoRows)
		mockRepo.AssertNotCalled(t, "UpdateExtracurricular")
	})

	t.Run("not found error is propagated", func(t *testing.T) {
		t.Parallel()

		input := &models.UpdateExtracurricularInput{ID: id, Body: models.UpdateExtracurricularRequest{Goals: ptr("x")}}
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricularStudentID", mock.Anything, id).Return(studentID, nil)
		mockRepo.On("UpdateExtracurricular", mock.Anything, id, *input).Return(nil, pgx.ErrNoRows)

		res, err := NewHandler(mockRepo, nil).UpdateExtracurricular(studentCtx(studentID, userID), id, input)

		assert.Nil(t, res)
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		t.Parallel()

		input := &models.UpdateExtracurricularInput{ID: id, Body: models.UpdateExtracurricularRequest{Goals: ptr("x")}}
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricularStudentID", mock.Anything, id).Return(studentID, nil)
		mockRepo.On("UpdateExtracurricular", mock.Anything, id, *input).Return(nil, errors.New("database error"))

		res, err := NewHandler(mockRepo, nil).UpdateExtracurricular(studentCtx(studentID, userID), id, input)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
