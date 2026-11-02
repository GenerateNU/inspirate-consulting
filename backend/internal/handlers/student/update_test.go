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

// sets the caller's role the way AuthMiddleware does
func withRole(ctx context.Context, role models.Role) context.Context {
	return context.WithValue(ctx, "Role", string(role)) //nolint:staticcheck
}

func TestHandler_UpdateStudent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("passes the update through and returns the student", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()
		gpa := 4
		input := models.UpdateStudentBody{GPA: &gpa}

		mockStudentRepo := mocks.NewStudentRepository(t)
		mockStudentRepo.On("UpdateStudent", ctx, studentID, input).
			Return(&models.Student{ID: studentID, GPA: gpa}, nil)

		handler := NewHandler(mocks.NewEssayReviewRepository(t), mockStudentRepo, nil)
		res, err := handler.UpdateStudent(ctx, studentID, input)

		require.NoError(t, err)
		assert.Equal(t, studentID, res.ID)
		assert.Equal(t, gpa, res.GPA)
	})

	t.Run("unknown student propagates the repository error", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()
		input := models.UpdateStudentBody{}

		mockStudentRepo := mocks.NewStudentRepository(t)
		mockStudentRepo.On("UpdateStudent", ctx, studentID, input).
			Return(nil, errs.NotFound("student", "id", studentID.String()))

		handler := NewHandler(mocks.NewEssayReviewRepository(t), mockStudentRepo, nil)
		res, err := handler.UpdateStudent(ctx, studentID, input)

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 404, err.(errs.HTTPError).Code)
	})

	t.Run("a student can update their own profile", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()
		gpa := 4
		input := models.UpdateStudentBody{GPA: &gpa}
		studentCtx := auth.WithStudentID(withRole(ctx, models.StudentRole), studentID)

		mockStudentRepo := mocks.NewStudentRepository(t)
		mockStudentRepo.On("UpdateStudent", studentCtx, studentID, input).
			Return(&models.Student{ID: studentID, GPA: gpa}, nil)

		handler := NewHandler(mocks.NewEssayReviewRepository(t), mockStudentRepo, auth.NewAuthorizer(nil))
		res, err := handler.UpdateStudent(studentCtx, studentID, input)

		require.NoError(t, err)
		assert.Equal(t, studentID, res.ID)
	})

	t.Run("a student cannot update another student", func(t *testing.T) {
		t.Parallel()

		studentCtx := auth.WithStudentID(withRole(ctx, models.StudentRole), uuid.New())

		mockStudentRepo := mocks.NewStudentRepository(t) // UpdateStudent should not be called

		handler := NewHandler(mocks.NewEssayReviewRepository(t), mockStudentRepo, auth.NewAuthorizer(nil))
		res, err := handler.UpdateStudent(studentCtx, uuid.New(), models.UpdateStudentBody{})

		assert.Nil(t, res)
		assert.Equal(t, errs.Forbidden("you can only access your own data"), err)
	})

	t.Run("a caller with an unknown role is rejected", func(t *testing.T) {
		t.Parallel()

		adminCtx := withRole(ctx, models.Role("admin"))

		mockStudentRepo := mocks.NewStudentRepository(t) // UpdateStudent should not be called

		handler := NewHandler(mocks.NewEssayReviewRepository(t), mockStudentRepo, auth.NewAuthorizer(nil))
		res, err := handler.UpdateStudent(adminCtx, uuid.New(), models.UpdateStudentBody{})

		assert.Nil(t, res)
		assert.Equal(t, errs.Forbidden("insufficient role"), err)
	})
}
