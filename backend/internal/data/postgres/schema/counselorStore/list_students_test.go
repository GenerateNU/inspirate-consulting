package counselorRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test list function returns only the counselor's students
func TestListStudents(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewCounselorRepository(db)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	otherCounselorID := createTestCounselor(t, db)

	alice := createTestStudent(t, db, counselorID, "Alice", "junior")
	carol := createTestStudent(t, db, counselorID, "Carol", "freshman")
	bob := createTestStudent(t, db, otherCounselorID, "Bob", "senior")

	results, err := repo.ListStudents(ctx, counselorID.String())
	if err != nil {
		t.Fatalf("ListStudents failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 students, got %d", len(results))
	}

	found := map[uuid.UUID]models.CounselorStudentRow{}
	for _, r := range results {
		found[r.StudentID] = r
	}
	if _, ok := found[bob]; ok {
		t.Error("expected the other counselor's student to NOT be present in results")
	}
	if _, ok := found[carol]; !ok {
		t.Error("expected Carol to be present in results")
	}

	got, ok := found[alice]
	if !ok {
		t.Fatal("expected Alice to be present in results")
	}
	if got.Name != "Alice" {
		t.Errorf("expected Name %q, got %q", "Alice", got.Name)
	}
	if got.Year != models.Junior {
		t.Errorf("expected Year %q, got %q", models.Junior, got.Year)
	}
}

func TestListStudents_NoStudents(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewCounselorRepository(db)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)

	results, err := repo.ListStudents(ctx, counselorID.String())
	if err != nil {
		t.Fatalf("ListStudents failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 students, got %d", len(results))
	}
}
