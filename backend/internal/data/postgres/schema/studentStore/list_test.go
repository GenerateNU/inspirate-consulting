package studentstore

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"
)

// seedStudents inserts students out of alphabetical order so the tests check the query's sort.
func seedStudents(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	_, counselorID := createTestCounselor(t, db, "List Counselor")
	for _, name := range []string{"Eva Park", "Anna Lee", "Joanna Smith", "Cara Diaz", "Hannah Kim", "Ben Cho", "Dan Wu"} {
		createTestStudent(t, db, name, counselorID)
	}
}

func listedNames(rows []models.StudentListItem) []string {
	names := make([]string, len(rows))
	for i, row := range rows {
		names[i] = row.Name
	}
	return names
}

func equalNames(t *testing.T, want, got []string) {
	t.Helper()

	if len(want) != len(got) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

func lastKey(rows []models.StudentListItem) *pagination.NameIDKey {
	last := rows[len(rows)-1]
	return &pagination.NameIDKey{Name: last.Name, ID: last.ID}
}

func TestListStudents_NoFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	seedStudents(t, db)

	// The handler sends an empty name when the name query param is omitted.
	empty := ""
	rows, err := repo.ListStudents(ctx, &empty, nil, 5)
	if err != nil {
		t.Fatalf("ListStudents failed: %v", err)
	}

	equalNames(t, []string{"Anna Lee", "Ben Cho", "Cara Diaz", "Dan Wu", "Eva Park"}, listedNames(rows))
}

func TestListStudents_CustomLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	seedStudents(t, db)

	rows, err := repo.ListStudents(ctx, nil, nil, 2)
	if err != nil {
		t.Fatalf("ListStudents failed: %v", err)
	}

	equalNames(t, []string{"Anna Lee", "Ben Cho"}, listedNames(rows))
}

func TestListStudents_ByName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	seedStudents(t, db)

	name := "ANN"
	rows, err := repo.ListStudents(ctx, &name, nil, 5)
	if err != nil {
		t.Fatalf("ListStudents failed: %v", err)
	}

	equalNames(t, []string{"Anna Lee", "Hannah Kim", "Joanna Smith"}, listedNames(rows))
}

func TestListStudents_NextCursor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	seedStudents(t, db)

	page1, err := repo.ListStudents(ctx, nil, nil, 3)
	if err != nil {
		t.Fatalf("ListStudents page 1 failed: %v", err)
	}
	equalNames(t, []string{"Anna Lee", "Ben Cho", "Cara Diaz"}, listedNames(page1))

	page2, err := repo.ListStudents(ctx, nil, lastKey(page1), 3)
	if err != nil {
		t.Fatalf("ListStudents page 2 failed: %v", err)
	}
	equalNames(t, []string{"Dan Wu", "Eva Park", "Hannah Kim"}, listedNames(page2))
}

func TestListStudents_LastPage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	seedStudents(t, db)

	page1, err := repo.ListStudents(ctx, nil, nil, 6)
	if err != nil {
		t.Fatalf("ListStudents page 1 failed: %v", err)
	}

	page2, err := repo.ListStudents(ctx, nil, lastKey(page1), 6)
	if err != nil {
		t.Fatalf("ListStudents page 2 failed: %v", err)
	}
	if len(page2) >= 6 {
		t.Fatalf("expected fewer rows than the limit on the last page, got %d", len(page2))
	}
	equalNames(t, []string{"Joanna Smith"}, listedNames(page2))
}
