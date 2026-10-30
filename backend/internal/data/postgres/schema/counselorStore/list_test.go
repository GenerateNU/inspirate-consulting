package counselorRepository

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"
)

// seedCounselors inserts counselors out of alphabetical order so the tests check the query's sort.
func seedCounselors(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	for _, name := range []string{"Sam Rivera", "David Chen", "Uma Reddy", "Maria Lopez", "Tara Singh", "James Okafor", "Priya Patel"} {
		createNamedTestCounselor(t, db, name)
	}
}

func listedNames(rows []models.CounselorListItem) []string {
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

func lastKey(rows []models.CounselorListItem) *pagination.NameIDKey {
	last := rows[len(rows)-1]
	return &pagination.NameIDKey{Name: last.Name, ID: last.ID}
}

func TestListCounselors_NoFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	seedCounselors(t, db)

	// The handler sends an empty name when the name query param is omitted.
	empty := ""
	rows, err := repo.ListCounselors(ctx, &empty, nil, 5)
	if err != nil {
		t.Fatalf("ListCounselors failed: %v", err)
	}

	equalNames(t, []string{"David Chen", "James Okafor", "Maria Lopez", "Priya Patel", "Sam Rivera"}, listedNames(rows))
}

func TestListCounselors_CustomLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	seedCounselors(t, db)

	rows, err := repo.ListCounselors(ctx, nil, nil, 2)
	if err != nil {
		t.Fatalf("ListCounselors failed: %v", err)
	}

	equalNames(t, []string{"David Chen", "James Okafor"}, listedNames(rows))
}

func TestListCounselors_ByName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	seedCounselors(t, db)

	name := "PAT"
	rows, err := repo.ListCounselors(ctx, &name, nil, 5)
	if err != nil {
		t.Fatalf("ListCounselors failed: %v", err)
	}

	equalNames(t, []string{"Priya Patel"}, listedNames(rows))
}

func TestListCounselors_NextCursor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	seedCounselors(t, db)

	page1, err := repo.ListCounselors(ctx, nil, nil, 3)
	if err != nil {
		t.Fatalf("ListCounselors page 1 failed: %v", err)
	}
	equalNames(t, []string{"David Chen", "James Okafor", "Maria Lopez"}, listedNames(page1))

	page2, err := repo.ListCounselors(ctx, nil, lastKey(page1), 3)
	if err != nil {
		t.Fatalf("ListCounselors page 2 failed: %v", err)
	}
	equalNames(t, []string{"Priya Patel", "Sam Rivera", "Tara Singh"}, listedNames(page2))
}

func TestListCounselors_LastPage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()
	seedCounselors(t, db)

	page1, err := repo.ListCounselors(ctx, nil, nil, 6)
	if err != nil {
		t.Fatalf("ListCounselors page 1 failed: %v", err)
	}

	page2, err := repo.ListCounselors(ctx, nil, lastKey(page1), 6)
	if err != nil {
		t.Fatalf("ListCounselors page 2 failed: %v", err)
	}
	if len(page2) >= 6 {
		t.Fatalf("expected fewer rows than the limit on the last page, got %d", len(page2))
	}
	equalNames(t, []string{"Uma Reddy"}, listedNames(page2))
}
