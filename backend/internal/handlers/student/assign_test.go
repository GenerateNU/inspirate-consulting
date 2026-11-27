package student

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func assignInput(userID uuid.UUID, counselorID *uuid.UUID) *models.AssignCounselorInput {
	input := &models.AssignCounselorInput{}
	input.Body.UserID = userID
	input.Body.CounselorID = counselorID
	return input
}

// expectStudentLookup mocks the handler's lookup of the student and its current
// counselor (nil when unassigned).
func expectStudentLookup(mockRepo *mocks.StudentRepository, ctx context.Context, userID, studentID uuid.UUID, current *uuid.UUID) {
	mockRepo.On("GetStudentIDByUserID", ctx, userID).Return(studentID, nil)
	mockRepo.On("GetStudent", ctx, studentID).
		Return(&models.Student{ID: studentID, UserID: userID, GraduationYear: ptr(2028), CounselorID: current}, nil)
}

// denyingAuthorizer is a real Authorizer with no database. With no role in the
// context it rejects every access check, so a test using it proves whether
// the handler ran the check.
func denyingAuthorizer() *auth.Authorizer {
	return auth.NewAuthorizer(nil)
}

func requireStatus(t *testing.T, err error, code int) {
	t.Helper()
	require.Error(t, err)
	var httpErr errs.HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, code, httpErr.Code)
}

// Unit tests for the AssignCounselor handler logic without HTTP or database dependencies.
// Mocks are strict: a repository call without an expectation fails the test.
func TestHandler_AssignCounselor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("any counselor can assign an unassigned student", func(t *testing.T) {
		t.Parallel()

		studentUserID, studentID, counselorID := uuid.New(), uuid.New(), uuid.New()
		updated := &models.Student{ID: studentID, UserID: studentUserID, GraduationYear: ptr(2028), CounselorID: &counselorID}

		mockRepo := mocks.NewStudentRepository(t)
		expectStudentLookup(mockRepo, ctx, studentUserID, studentID, nil)
		mockRepo.On("AssignCounselor", ctx, studentUserID, &counselorID).Return(updated, nil)

		// The authorizer would deny, so success proves no access check ran.
		handler := NewHandler(nil, mockRepo, denyingAuthorizer())
		res, err := handler.AssignCounselor(ctx, assignInput(studentUserID, &counselorID))

		require.NoError(t, err)
		assert.Equal(t, updated, res)
	})

	t.Run("the current counselor can reassign their student", func(t *testing.T) {
		t.Parallel()

		studentUserID, studentID := uuid.New(), uuid.New()
		currentID, newID := uuid.New(), uuid.New()
		updated := &models.Student{ID: studentID, UserID: studentUserID, GraduationYear: ptr(2028), CounselorID: &newID}

		mockRepo := mocks.NewStudentRepository(t)
		expectStudentLookup(mockRepo, ctx, studentUserID, studentID, &currentID)
		mockRepo.On("AssignCounselor", ctx, studentUserID, &newID).Return(updated, nil)

		// A nil authorizer (test mode) allows access.
		handler := NewHandler(nil, mockRepo, nil)
		res, err := handler.AssignCounselor(ctx, assignInput(studentUserID, &newID))

		require.NoError(t, err)
		require.NotNil(t, res.CounselorID)
		assert.Equal(t, newID, *res.CounselorID)
	})

	t.Run("reassigning another counselor's student is forbidden", func(t *testing.T) {
		t.Parallel()

		studentUserID, studentID := uuid.New(), uuid.New()
		currentID, newID := uuid.New(), uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		expectStudentLookup(mockRepo, ctx, studentUserID, studentID, &currentID)
		// AssignCounselor must not be called: the access check runs before the write.

		handler := NewHandler(nil, mockRepo, denyingAuthorizer())
		res, err := handler.AssignCounselor(ctx, assignInput(studentUserID, &newID))

		assert.Nil(t, res)
		requireStatus(t, err, 403)
	})

	t.Run("the current counselor can unassign their student", func(t *testing.T) {
		t.Parallel()

		studentUserID, studentID, currentID := uuid.New(), uuid.New(), uuid.New()
		updated := &models.Student{ID: studentID, UserID: studentUserID, GraduationYear: ptr(2028)}

		mockRepo := mocks.NewStudentRepository(t)
		expectStudentLookup(mockRepo, ctx, studentUserID, studentID, &currentID)
		mockRepo.On("AssignCounselor", ctx, studentUserID, (*uuid.UUID)(nil)).Return(updated, nil)

		handler := NewHandler(nil, mockRepo, nil)
		res, err := handler.AssignCounselor(ctx, assignInput(studentUserID, nil))

		require.NoError(t, err)
		assert.Nil(t, res.CounselorID)
	})

	t.Run("unassigning another counselor's student is forbidden", func(t *testing.T) {
		t.Parallel()

		studentUserID, studentID, currentID := uuid.New(), uuid.New(), uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		expectStudentLookup(mockRepo, ctx, studentUserID, studentID, &currentID)
		// AssignCounselor must not be called.

		handler := NewHandler(nil, mockRepo, denyingAuthorizer())
		res, err := handler.AssignCounselor(ctx, assignInput(studentUserID, nil))

		assert.Nil(t, res)
		requireStatus(t, err, 403)
	})

	t.Run("unknown student returns not found", func(t *testing.T) {
		t.Parallel()

		studentUserID, counselorID := uuid.New(), uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("GetStudentIDByUserID", ctx, studentUserID).
			Return(uuid.Nil, errs.NotFound("student", "user_id", studentUserID))

		handler := NewHandler(nil, mockRepo, nil)
		res, err := handler.AssignCounselor(ctx, assignInput(studentUserID, &counselorID))

		assert.Nil(t, res)
		requireStatus(t, err, 404)
	})

	t.Run("unassigning an unknown student returns not found", func(t *testing.T) {
		t.Parallel()

		studentUserID := uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("GetStudentIDByUserID", ctx, studentUserID).
			Return(uuid.Nil, errs.NotFound("student", "user_id", studentUserID))

		handler := NewHandler(nil, mockRepo, nil)
		res, err := handler.AssignCounselor(ctx, assignInput(studentUserID, nil))

		assert.Nil(t, res)
		requireStatus(t, err, 404)
	})

	t.Run("unknown counselor returns not found", func(t *testing.T) {
		t.Parallel()

		studentUserID, studentID, counselorID := uuid.New(), uuid.New(), uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		expectStudentLookup(mockRepo, ctx, studentUserID, studentID, nil)
		mockRepo.On("AssignCounselor", ctx, studentUserID, &counselorID).
			Return(nil, errs.NotFound("counselor", "id", counselorID))

		handler := NewHandler(nil, mockRepo, nil)
		res, err := handler.AssignCounselor(ctx, assignInput(studentUserID, &counselorID))

		assert.Nil(t, res)
		requireStatus(t, err, 404)
	})

	t.Run("assigning a counselor to another counselor fails", func(t *testing.T) {
		t.Parallel()

		// A counselor's user id has no student row, so the lookup finds no student.
		counselorUserID, counselorID := uuid.New(), uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("GetStudentIDByUserID", ctx, counselorUserID).
			Return(uuid.Nil, errs.NotFound("student", "user_id", counselorUserID))

		handler := NewHandler(nil, mockRepo, nil)
		res, err := handler.AssignCounselor(ctx, assignInput(counselorUserID, &counselorID))

		assert.Nil(t, res)
		requireStatus(t, err, 404)
	})
}
