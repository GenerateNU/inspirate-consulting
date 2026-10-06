package mediaAccessRepository

import (
	"context"
	"testing"

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

	results, err := repo.ListAccessibleMedia(ctx, studentID, 10, 0)
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
			results, err := repo.ListAccessibleMedia(ctx, studentID, tc.limit, tc.offset)
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

	results, err := repo.ListAccessibleMedia(ctx, uuid.NewString(), 10, 0)
	if err != nil {
		t.Fatalf("ListAccessibleMedia failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 media for a student with no access, got %d", len(results))
	}
}
