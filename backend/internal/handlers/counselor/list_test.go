package counselor

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

func counselorRows(names ...string) []models.CounselorListItem {
	rows := make([]models.CounselorListItem, len(names))
	for i, name := range names {
		rows[i] = models.CounselorListItem{ID: uuid.New(), UserID: uuid.New(), Name: name}
	}
	return rows
}

func counselorNames(rows []models.CounselorListItem) []string {
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

// Unit tests for the ListCounselors handler logic without HTTP or database dependencies.
func TestHandler_ListCounselors(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("no filters returns the first page and a next cursor", func(t *testing.T) {
		t.Parallel()

		rows := counselorRows("David Chen", "James Okafor", "Maria Lopez", "Priya Patel", "Sam Rivera", "Tara Singh")

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListCounselors", ctx, nameIs(""), noCursor(), 6).Return(rows, nil)

		handler := NewHandler(mockRepo)
		page, err := handler.ListCounselors(ctx, &models.FetchCounselorsInput{Limit: 5})

		require.NoError(t, err)
		assert.Equal(t, []string{"David Chen", "James Okafor", "Maria Lopez", "Priya Patel", "Sam Rivera"}, counselorNames(page.Counselors))
		require.NotNil(t, page.NextCursor)

		key, err := pagination.Decode[pagination.NameIDKey](*page.NextCursor)
		require.NoError(t, err)
		assert.Equal(t, pagination.NameIDKey{Name: "Sam Rivera", ID: rows[4].ID}, key)
	})

	t.Run("a custom limit asks for one extra row and returns limit rows", func(t *testing.T) {
		t.Parallel()

		rows := counselorRows("David Chen", "James Okafor", "Maria Lopez")

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListCounselors", ctx, nameIs(""), noCursor(), 3).Return(rows, nil)

		handler := NewHandler(mockRepo)
		page, err := handler.ListCounselors(ctx, &models.FetchCounselorsInput{Limit: 2})

		require.NoError(t, err)
		assert.Equal(t, []string{"David Chen", "James Okafor"}, counselorNames(page.Counselors))
		assert.NotNil(t, page.NextCursor)
	})

	t.Run("a name is passed to the repository as the search term", func(t *testing.T) {
		t.Parallel()

		rows := counselorRows("Priya Patel")

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListCounselors", ctx, nameIs("pat"), noCursor(), 6).Return(rows, nil)

		handler := NewHandler(mockRepo)
		page, err := handler.ListCounselors(ctx, &models.FetchCounselorsInput{Name: "pat", Limit: 5})

		require.NoError(t, err)
		assert.Equal(t, []string{"Priya Patel"}, counselorNames(page.Counselors))
		assert.Nil(t, page.NextCursor)
	})

	t.Run("the next cursor continues after the last row of the previous page", func(t *testing.T) {
		t.Parallel()

		first := counselorRows("David Chen", "James Okafor", "Maria Lopez")
		second := counselorRows("Maria Lopez", "Priya Patel")

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListCounselors", ctx, nameIs(""), noCursor(), 3).Return(first, nil).Once()
		mockRepo.On("ListCounselors", ctx, nameIs(""), cursorAt("James Okafor", first[1].ID), 3).Return(second, nil).Once()

		handler := NewHandler(mockRepo)

		page1, err := handler.ListCounselors(ctx, &models.FetchCounselorsInput{Limit: 2})
		require.NoError(t, err)
		assert.Equal(t, []string{"David Chen", "James Okafor"}, counselorNames(page1.Counselors))
		require.NotNil(t, page1.NextCursor)

		page2, err := handler.ListCounselors(ctx, &models.FetchCounselorsInput{Limit: 2, Cursor: *page1.NextCursor})
		require.NoError(t, err)
		assert.Equal(t, []string{"Maria Lopez", "Priya Patel"}, counselorNames(page2.Counselors))
	})

	t.Run("fewer rows than the limit is the last page", func(t *testing.T) {
		t.Parallel()

		rows := counselorRows("Priya Patel", "Sam Rivera")

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListCounselors", ctx, nameIs(""), noCursor(), 6).Return(rows, nil)

		handler := NewHandler(mockRepo)
		page, err := handler.ListCounselors(ctx, &models.FetchCounselorsInput{Limit: 5})

		require.NoError(t, err)
		assert.Len(t, page.Counselors, 2)
		assert.Nil(t, page.NextCursor)
	})
}
