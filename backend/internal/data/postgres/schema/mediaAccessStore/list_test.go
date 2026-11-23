package mediaAccessRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test list function returns only the media a student has been granted access to
func TestListAccessibleMedia(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	studentID := uuid.NewString()
	otherStudentID := uuid.NewString()

	mediaA := createTestMedia(t, db, "List Test Alpha")
	mediaB := createTestMedia(t, db, "List Test Beta")
	otherMedia := createTestMedia(t, db, "List Test Other")

	grants := []*models.GrantMediaAccessRequestBody{
		{StudentID: studentID, MediaID: mediaA},
		{StudentID: studentID, MediaID: mediaB},
		// A different student's access should not show up in studentID's results.
		{StudentID: otherStudentID, MediaID: otherMedia},
	}
	for _, g := range grants {
		if _, err := repo.GrantMediaAccess(ctx, g); err != nil {
			t.Fatalf("setup GrantMediaAccess failed: %v", err)
		}
	}

	results, err := repo.ListAccessibleMedia(ctx, studentID, "", 0.2, 10, 0)
	if err != nil {
		t.Fatalf("ListAccessibleMedia failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 media for studentID, got %d", len(results))
	}

	foundIDs := map[string]bool{}
	for _, r := range results {
		foundIDs[r.ID] = true
	}
	if !foundIDs[mediaA] || !foundIDs[mediaB] {
		t.Error("expected both of studentID's media to be present in results")
	}
	if foundIDs[otherMedia] {
		t.Error("expected the other student's media to NOT be present in results")
	}
}

// test list function respects limit and offset, returning results ordered by title
func TestListAccessibleMedia_Pagination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	studentID := uuid.NewString()

	titles := []string{"Page Test A", "Page Test B", "Page Test C"}
	for _, title := range titles {
		mediaID := createTestMedia(t, db, title)
		if _, err := repo.GrantMediaAccess(ctx, &models.GrantMediaAccessRequestBody{
			StudentID: studentID,
			MediaID:   mediaID,
		}); err != nil {
			t.Fatalf("setup GrantMediaAccess failed: %v", err)
		}
	}

	tests := []struct {
		name           string
		limit          int
		offset         int
		expectedTitles []string
	}{
		{name: "first page", limit: 2, offset: 0, expectedTitles: []string{"Page Test A", "Page Test B"}},
		{name: "second page", limit: 2, offset: 2, expectedTitles: []string{"Page Test C"}},
		{name: "offset past end", limit: 2, offset: 3, expectedTitles: []string{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results, err := repo.ListAccessibleMedia(ctx, studentID, "", 0.2, tc.limit, tc.offset)
			if err != nil {
				t.Fatalf("ListAccessibleMedia failed: %v", err)
			}

			if len(results) != len(tc.expectedTitles) {
				t.Fatalf("expected %d media, got %d", len(tc.expectedTitles), len(results))
			}
			for i, r := range results {
				if r.Title != tc.expectedTitles[i] {
					t.Errorf("result %d: expected title %q, got %q", i, tc.expectedTitles[i], r.Title)
				}
			}
		})
	}
}

// test list function returns an empty result, not an error, for a student with no access
func TestListAccessibleMedia_NoAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo := NewMediaAccessRepository(testutils.SetupTestDB(t))
	ctx := context.Background()

	results, err := repo.ListAccessibleMedia(ctx, uuid.NewString(), "", 0.2, 10, 0)
	if err != nil {
		t.Fatalf("ListAccessibleMedia failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 media for a student with no access, got %d", len(results))
	}
}

// test each listed media carries the student's own progress and their todo's deadline
func TestListAccessibleMedia_ProgressAndDueDate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	studentID := uuid.NewString()
	otherStudentID := uuid.NewString()
	withTodos := createTestMedia(t, db, "Due Date A")
	noTodos := createTestMedia(t, db, "Due Date B")
	for _, g := range []*models.GrantMediaAccessRequestBody{
		{StudentID: studentID, MediaID: withTodos},
		{StudentID: studentID, MediaID: noTodos},
		{StudentID: otherStudentID, MediaID: withTodos},
		{StudentID: otherStudentID, MediaID: noTodos},
	} {
		if _, err := repo.GrantMediaAccess(ctx, g); err != nil {
			t.Fatalf("setup GrantMediaAccess failed: %v", err)
		}
	}

	// Assigned twice: the most recently created todo's deadline wins, even once that todo is
	// completed (the due date is the todo's, not "next open deadline").
	firstAssigned := time.Date(2026, 12, 1, 23, 59, 0, 0, time.UTC)
	latestAssigned := time.Date(2026, 12, 5, 23, 59, 0, 0, time.UTC)
	createTestMediaTodo(t, db, studentID, withTodos, firstAssigned, false)
	createTestMediaTodo(t, db, studentID, withTodos, latestAssigned, true)
	// another student's todo for the same media, created last, must not count
	createTestMediaTodo(t, db, otherStudentID, withTodos, latestAssigned.Add(48*time.Hour), false)

	if _, err := repo.RecordMediaProgress(ctx, studentID, withTodos, true); err != nil {
		t.Fatalf("setup RecordMediaProgress failed: %v", err)
	}
	// the other student's progress must not leak into this student's list
	if _, err := repo.RecordMediaProgress(ctx, otherStudentID, noTodos, true); err != nil {
		t.Fatalf("setup RecordMediaProgress failed: %v", err)
	}

	results, err := repo.ListAccessibleMedia(ctx, studentID, "", 0.2, 10, 0)
	if err != nil {
		t.Fatalf("ListAccessibleMedia failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 media, got %d", len(results))
	}

	a, b := results[0], results[1] // ordered by title
	if a.ID != withTodos || b.ID != noTodos {
		t.Fatalf("unexpected order: %s, %s", a.ID, b.ID)
	}
	if a.DueDate == nil || !a.DueDate.Equal(latestAssigned) {
		t.Errorf("expected DueDate %v, got %v", latestAssigned, a.DueDate)
	}
	if a.StartedAt == nil || a.CompletedAt == nil {
		t.Errorf("expected progress on media A, got started=%v completed=%v", a.StartedAt, a.CompletedAt)
	}
	if b.DueDate != nil || b.StartedAt != nil || b.CompletedAt != nil {
		t.Errorf("expected no due date or progress on media B, got %+v", b)
	}
}

// test the search query filters by title or description, case-insensitively
func TestListAccessibleMedia_Search(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewMediaAccessRepository(db)
	ctx := context.Background()

	studentID := uuid.NewString()
	match := createTestMedia(t, db, "FAFSA Basics")
	noMatch := createTestMedia(t, db, "Essay Tips")
	for _, mediaID := range []string{match, noMatch} {
		if _, err := repo.GrantMediaAccess(ctx, &models.GrantMediaAccessRequestBody{StudentID: studentID, MediaID: mediaID}); err != nil {
			t.Fatalf("setup GrantMediaAccess failed: %v", err)
		}
	}

	results, err := repo.ListAccessibleMedia(ctx, studentID, "fafsa", 0.4, 10, 0)
	if err != nil {
		t.Fatalf("ListAccessibleMedia failed: %v", err)
	}
	if len(results) != 1 || results[0].ID != match {
		t.Errorf("expected only %s, got %+v", match, results)
	}
}
