package student

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"
)

func studentRows(names ...string) []models.StudentListItem {
	rows := make([]models.StudentListItem, len(names))
	for i, name := range names {
		rows[i] = models.StudentListItem{ID: uuid.New(), UserID: uuid.New(), Name: name, Year: models.Junior}
	}
	return rows
}

func studentNames(rows []models.StudentListItem) []string {
	names := make([]string, len(rows))
	for i, row := range rows {
		names[i] = row.Name
	}
	return names
}

func nameIs(want string) any {
	return mock.MatchedBy(func(name *string) bool { return name != nil && *name == want })
}

func noCursor() any {
	return mock.MatchedBy(func(after *pagination.NameIDKey) bool { return after == nil })
}

func cursorAt(name string, id uuid.UUID) any {
	return mock.MatchedBy(func(after *pagination.NameIDKey) bool {
		return after != nil && after.Name == name && after.ID == id
	})
}

// Unit tests for the ListStudents handler logic without HTTP or database dependencies.
func TestHandler_ListStudents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("no filters returns the first page and a next cursor", func(t *testing.T) {
		t.Parallel()

		rows := studentRows("Anna Lee", "Ben Cho", "Cara Diaz", "Dan Wu", "Eva Park", "Finn Ross")

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("ListStudents", ctx, nameIs(""), noCursor(), 6).Return(rows, nil)

		handler := NewHandler(nil, mockRepo, nil)
		page, err := handler.ListStudents(ctx, &models.FetchStudentsInput{Limit: 5})

		require.NoError(t, err)
		assert.Equal(t, []string{"Anna Lee", "Ben Cho", "Cara Diaz", "Dan Wu", "Eva Park"}, studentNames(page.Students))
		require.NotNil(t, page.NextCursor)

		key, err := pagination.Decode[pagination.NameIDKey](*page.NextCursor)
		require.NoError(t, err)
		assert.Equal(t, pagination.NameIDKey{Name: "Eva Park", ID: rows[4].ID}, key)
	})

	t.Run("a custom limit asks for one extra row and returns limit rows", func(t *testing.T) {
		t.Parallel()

		rows := studentRows("Anna Lee", "Ben Cho", "Cara Diaz", "Dan Wu")

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("ListStudents", ctx, nameIs(""), noCursor(), 4).Return(rows, nil)

		handler := NewHandler(nil, mockRepo, nil)
		page, err := handler.ListStudents(ctx, &models.FetchStudentsInput{Limit: 3})

		require.NoError(t, err)
		assert.Equal(t, []string{"Anna Lee", "Ben Cho", "Cara Diaz"}, studentNames(page.Students))
		assert.NotNil(t, page.NextCursor)
	})

	t.Run("a name is passed to the repository as the search term", func(t *testing.T) {
		t.Parallel()

		rows := studentRows("Anna Lee", "Hannah Kim")

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("ListStudents", ctx, nameIs("ann"), noCursor(), 6).Return(rows, nil)

		handler := NewHandler(nil, mockRepo, nil)
		page, err := handler.ListStudents(ctx, &models.FetchStudentsInput{Name: "ann", Limit: 5})

		require.NoError(t, err)
		assert.Equal(t, []string{"Anna Lee", "Hannah Kim"}, studentNames(page.Students))
		assert.Nil(t, page.NextCursor)
	})

	t.Run("the next cursor continues after the last row of the previous page", func(t *testing.T) {
		t.Parallel()

		first := studentRows("Anna Lee", "Ben Cho", "Cara Diaz")
		second := studentRows("Dan Wu", "Eva Park")

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("ListStudents", ctx, nameIs(""), noCursor(), 3).Return(first, nil).Once()
		mockRepo.On("ListStudents", ctx, nameIs(""), cursorAt("Ben Cho", first[1].ID), 3).Return(second, nil).Once()

		handler := NewHandler(nil, mockRepo, nil)

		page1, err := handler.ListStudents(ctx, &models.FetchStudentsInput{Limit: 2})
		require.NoError(t, err)
		assert.Equal(t, []string{"Anna Lee", "Ben Cho"}, studentNames(page1.Students))
		require.NotNil(t, page1.NextCursor)

		page2, err := handler.ListStudents(ctx, &models.FetchStudentsInput{Limit: 2, Cursor: *page1.NextCursor})
		require.NoError(t, err)
		assert.Equal(t, []string{"Dan Wu", "Eva Park"}, studentNames(page2.Students))
	})

	t.Run("fewer rows than the limit is the last page", func(t *testing.T) {
		t.Parallel()

		rows := studentRows("Dan Wu", "Eva Park", "Finn Ross")

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("ListStudents", ctx, nameIs(""), noCursor(), 6).Return(rows, nil)

		handler := NewHandler(nil, mockRepo, nil)
		page, err := handler.ListStudents(ctx, &models.FetchStudentsInput{Limit: 5})

		require.NoError(t, err)
		assert.Len(t, page.Students, 3)
		assert.Nil(t, page.NextCursor)
	})
}
