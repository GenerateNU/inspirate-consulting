package counselor

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
)

var (
	aliceID = uuid.MustParse("00000000-0000-0000-0000-00000000000a")
	bobID   = uuid.MustParse("00000000-0000-0000-0000-00000000000b")
	carolID = uuid.MustParse("00000000-0000-0000-0000-00000000000c")
	daveID  = uuid.MustParse("00000000-0000-0000-0000-00000000000d")
)

// roster for testing
// Alice (junior): pending review
// Bob (senior): falling behind
// Carol (freshman): neither
// Dave (sophomore): pending review and falling behind
func testRoster() []models.CounselorStudentRow {
	return []models.CounselorStudentRow{
		{StudentID: daveID, Name: "Dave", Year: models.Sophomore},
		{StudentID: aliceID, Name: "Alice", Year: models.Junior},
		{StudentID: carolID, Name: "Carol", Year: models.Freshman},
		{StudentID: bobID, Name: "Bob", Year: models.Senior},
	}
}

// Alice has two open reviews, to check she is still flagged (and listed) once
func testReviews() []models.PendingReview {
	return []models.PendingReview{
		{TransactionID: uuid.New(), StudentID: aliceID, StudentName: "Alice"},
		{TransactionID: uuid.New(), StudentID: daveID, StudentName: "Dave"},
		{TransactionID: uuid.New(), StudentID: aliceID, StudentName: "Alice"},
	}
}

func testFallingBehindIDs() []uuid.UUID {
	return []uuid.UUID{bobID, daveID}
}

// sets up a mock returning the test roster, reviews, and falling behind ids
func newTestRepo(t *testing.T, ctx context.Context) *mocks.CounselorRepository {
	counselorID := auth.GetCounselorID(ctx)

	mockRepo := mocks.NewCounselorRepository(t)
	mockRepo.On("ListStudents", ctx, counselorID).Return(testRoster(), nil)
	mockRepo.On("ListPendingReviews", ctx, counselorID).Return(testReviews(), nil)
	mockRepo.On("ListFallingBehindStudentIDs", ctx, counselorID).Return(testFallingBehindIDs(), nil)
	return mockRepo
}

func names(students []models.CounselorStudent) []string {
	result := make([]string, 0, len(students))
	for _, student := range students {
		result = append(result, student.Name)
	}
	return result
}

// Unit tests for the ListStudents handler logic
// Written with stub counselor ID
func TestHandler_ListStudents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	counselorID := auth.GetCounselorID(ctx)

	t.Run("success sets flags", func(t *testing.T) {
		t.Parallel()

		handler := NewHandler(newTestRepo(t, ctx))
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{SortBy: "name", SortOrder: "asc"})

		assert.NoError(t, err)
		assert.Equal(t, []models.CounselorStudent{
			{CounselorStudentRow: models.CounselorStudentRow{StudentID: aliceID, Name: "Alice", Year: models.Junior}, IsFallingBehind: false, HasPendingReview: true},
			{CounselorStudentRow: models.CounselorStudentRow{StudentID: bobID, Name: "Bob", Year: models.Senior}, IsFallingBehind: true, HasPendingReview: false},
			{CounselorStudentRow: models.CounselorStudentRow{StudentID: carolID, Name: "Carol", Year: models.Freshman}, IsFallingBehind: false, HasPendingReview: false},
			{CounselorStudentRow: models.CounselorStudentRow{StudentID: daveID, Name: "Dave", Year: models.Sophomore}, IsFallingBehind: true, HasPendingReview: true},
		}, res)
	})

	t.Run("filters by year", func(t *testing.T) {
		t.Parallel()

		handler := NewHandler(newTestRepo(t, ctx))
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{Year: models.Junior, SortBy: "name", SortOrder: "asc"})

		assert.NoError(t, err)
		assert.Equal(t, []string{"Alice"}, names(res))
	})

	t.Run("filters by pending review", func(t *testing.T) {
		t.Parallel()

		handler := NewHandler(newTestRepo(t, ctx))
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{HasPendingReview: true, SortBy: "name", SortOrder: "asc"})

		assert.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Dave"}, names(res))
	})

	t.Run("combines filters", func(t *testing.T) {
		t.Parallel()

		handler := NewHandler(newTestRepo(t, ctx))
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{Year: models.Sophomore, HasPendingReview: true, SortBy: "name", SortOrder: "asc"})

		assert.NoError(t, err)
		assert.Equal(t, []string{"Dave"}, names(res))
	})

	t.Run("filters with no matches return an empty list", func(t *testing.T) {
		t.Parallel()

		handler := NewHandler(newTestRepo(t, ctx))
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{Year: models.Senior, HasPendingReview: true, SortBy: "name", SortOrder: "asc"})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res)
	})

	t.Run("sorts by name desc", func(t *testing.T) {
		t.Parallel()

		handler := NewHandler(newTestRepo(t, ctx))
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{SortBy: "name", SortOrder: "desc"})

		assert.NoError(t, err)
		assert.Equal(t, []string{"Dave", "Carol", "Bob", "Alice"}, names(res))
	})

	t.Run("sorts by year from freshman to senior, not alphabetically", func(t *testing.T) {
		t.Parallel()

		handler := NewHandler(newTestRepo(t, ctx))
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{SortBy: "year", SortOrder: "asc"})

		assert.NoError(t, err)
		assert.Equal(t, []string{"Carol", "Dave", "Alice", "Bob"}, names(res))
	})

	t.Run("sorts by year desc", func(t *testing.T) {
		t.Parallel()

		handler := NewHandler(newTestRepo(t, ctx))
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{SortBy: "year", SortOrder: "desc"})

		assert.NoError(t, err)
		assert.Equal(t, []string{"Bob", "Alice", "Dave", "Carol"}, names(res))
	})

	t.Run("year ties fall back to name, then student id", func(t *testing.T) {
		t.Parallel()

		sameNameID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		roster := []models.CounselorStudentRow{
			{StudentID: bobID, Name: "Zoe", Year: models.Junior},
			{StudentID: aliceID, Name: "Amy", Year: models.Junior},
			{StudentID: sameNameID, Name: "Zoe", Year: models.Junior},
		}

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListStudents", ctx, counselorID).Return(roster, nil)
		mockRepo.On("ListPendingReviews", ctx, counselorID).Return([]models.PendingReview{}, nil)
		mockRepo.On("ListFallingBehindStudentIDs", ctx, counselorID).Return([]uuid.UUID{}, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{SortBy: "year", SortOrder: "asc"})

		assert.NoError(t, err)
		assert.Equal(t, []uuid.UUID{aliceID, sameNameID, bobID}, []uuid.UUID{res[0].StudentID, res[1].StudentID, res[2].StudentID})
	})

	t.Run("no students returns an empty list", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListStudents", ctx, counselorID).Return([]models.CounselorStudentRow{}, nil)
		mockRepo.On("ListPendingReviews", ctx, counselorID).Return([]models.PendingReview{}, nil)
		mockRepo.On("ListFallingBehindStudentIDs", ctx, counselorID).Return([]uuid.UUID{}, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{SortBy: "name", SortOrder: "asc"})

		assert.NoError(t, err)
		assert.Equal(t, []models.CounselorStudent{}, res)
	})

	t.Run("list students error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListStudents", ctx, counselorID).Return(nil, errors.New("database error"))
		// ListPendingReviews and ListFallingBehindStudentIDs are not called if listing students fails

		handler := NewHandler(mockRepo)
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{SortBy: "name", SortOrder: "asc"})

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})

	t.Run("list pending reviews error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListStudents", ctx, counselorID).Return(testRoster(), nil)
		mockRepo.On("ListPendingReviews", ctx, counselorID).Return(nil, errors.New("database error"))
		// ListFallingBehindStudentIDs is not called if listing pending reviews fails

		handler := NewHandler(mockRepo)
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{SortBy: "name", SortOrder: "asc"})

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})

	t.Run("list falling behind students error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListStudents", ctx, counselorID).Return(testRoster(), nil)
		mockRepo.On("ListPendingReviews", ctx, counselorID).Return(testReviews(), nil)
		mockRepo.On("ListFallingBehindStudentIDs", ctx, counselorID).Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.ListStudents(ctx, models.ListCounselorStudentsInput{SortBy: "name", SortOrder: "asc"})

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
