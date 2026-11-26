package extracurricular

import (
	"context"
	"errors"
	"testing"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Unit tests for the ListExtracurriculars handler.
func TestHandler_ListExtracurriculars(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	userID := uuid.New()

	t.Run("success returns the student's rows including new fields", func(t *testing.T) {
		t.Parallel()

		rows := []models.Extracurricular{
			{ID: 1, StudentID: studentID.String(), Name: "Debate", HoursPerWeek: ptr(5.0), WeeksPerYear: ptr(int16(40)), ReadyForApplication: true, Goals: ptr("Nationals")},
			{ID: 2, StudentID: studentID.String(), Name: "Chess"},
		}
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("ListExtracurriculars", mock.Anything, studentID.String()).Return(rows, nil)

		res, err := NewHandler(mockRepo, nil).ListExtracurriculars(studentCtx(studentID, userID), nil)

		require.NoError(t, err)
		assert.Equal(t, rows, res)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("ListExtracurriculars", mock.Anything, studentID.String()).Return([]models.Extracurricular{}, nil)

		res, err := NewHandler(mockRepo, nil).ListExtracurriculars(studentCtx(studentID, userID), nil)

		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res)
	})

	t.Run("no student id in the request or the auth context is a bad request", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)

		res, err := NewHandler(mockRepo, nil).ListExtracurriculars(auth.WithUserID(context.Background(), userID), nil)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, 400, httpErr.Code)
		mockRepo.AssertNotCalled(t, "ListExtracurriculars")
	})

	t.Run("requested student id is listed when the caller may access it", func(t *testing.T) {
		t.Parallel()

		requested := uuid.New()
		rows := []models.Extracurricular{{ID: 3, StudentID: requested.String(), Name: "Robotics"}}
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("ListExtracurriculars", mock.Anything, requested.String()).Return(rows, nil)

		// no student in auth (a counselor); nil Authorizer allows access, so this
		// covers the student_id plumbing only
		raw := requested.String()
		res, err := NewHandler(mockRepo, nil).ListExtracurriculars(auth.WithUserID(context.Background(), userID), &raw)

		require.NoError(t, err)
		assert.Equal(t, rows, res)
	})

	t.Run("student always gets their own list, even when passing another student's id", func(t *testing.T) {
		t.Parallel()

		rows := []models.Extracurricular{{ID: 4, StudentID: studentID.String(), Name: "Choir"}}
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("ListExtracurriculars", mock.Anything, studentID.String()).Return(rows, nil)

		other := uuid.NewString()
		res, err := NewHandler(mockRepo, realAuthorizer()).ListExtracurriculars(studentRoleCtx(studentID), &other)

		require.NoError(t, err)
		assert.Equal(t, rows, res)
		mockRepo.AssertNotCalled(t, "ListExtracurriculars", mock.Anything, other)
	})

	t.Run("malformed student id is a bad request", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)

		raw := "not-a-uuid"
		res, err := NewHandler(mockRepo, nil).ListExtracurriculars(auth.WithUserID(context.Background(), userID), &raw)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, 400, httpErr.Code)
		mockRepo.AssertNotCalled(t, "ListExtracurriculars")
	})

	t.Run("repository error is propagated", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("ListExtracurriculars", mock.Anything, studentID.String()).Return(nil, errors.New("database error"))

		res, err := NewHandler(mockRepo, nil).ListExtracurriculars(studentCtx(studentID, userID), nil)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
