package extracurricularRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
)

func assertNotFound(t *testing.T, err error) {
	t.Helper()
	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != http.StatusNotFound {
		t.Errorf("expected 404 HTTPError, got %v", err)
	}
}

// Deleting removes the row and cascades to its achievements and notes.
func TestDeleteExtracurricular_CascadesToChildren(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	studentID := uuid.NewString()
	created := createTestExtracurricular(t, repo, studentID, createTestUser(t, db), fullCreateInput())
	insertAchievement(t, db, created.ID, "2024-05-01", "Award")
	insertAchievement(t, db, created.ID, "2024-06-01", "Another award")
	insertNote(t, db, created.ID, "2024-05-01T00:00:00Z", "Note")

	if err := repo.DeleteExtracurricular(context.Background(), created.ID, studentID); err != nil {
		t.Fatalf("DeleteExtracurricular failed: %v", err)
	}

	for _, table := range []string{"extracurriculars", "extracurricular_achievements", "extracurricular_notes"} {
		if n := countRows(t, db, table, created.ID); n != 0 {
			t.Errorf("%s still has %d rows for id %d", table, n, created.ID)
		}
	}

	_, err := repo.GetExtracurricular(context.Background(), created.ID)
	assertNotFound(t, err)
}

// Other extracurriculars and their children are left alone.
func TestDeleteExtracurricular_LeavesOtherRows(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	studentID := uuid.NewString()
	userID := createTestUser(t, db)
	target := createTestExtracurricular(t, repo, studentID, userID, minimalCreateInput("Target", "2025-01-01"))
	keep := createTestExtracurricular(t, repo, studentID, userID, minimalCreateInput("Keep", "2024-01-01"))
	insertAchievement(t, db, keep.ID, "2024-05-01", "Keep award")
	insertNote(t, db, keep.ID, "2024-05-01T00:00:00Z", "Keep note")

	if err := repo.DeleteExtracurricular(context.Background(), target.ID, studentID); err != nil {
		t.Fatalf("DeleteExtracurricular failed: %v", err)
	}

	got, err := repo.GetExtracurricular(context.Background(), keep.ID)
	if err != nil {
		t.Fatalf("kept row is gone: %v", err)
	}
	if len(got.Achievements) != 1 || len(got.Notes) != 1 {
		t.Errorf("kept row lost children: achievements=%d notes=%d", len(got.Achievements), len(got.Notes))
	}

	list, err := repo.ListExtracurriculars(context.Background(), studentID)
	if err != nil {
		t.Fatalf("ListExtracurriculars failed: %v", err)
	}
	if len(list) != 1 || list[0].ID != keep.ID {
		t.Errorf("expected only the kept row to remain, got %+v", list)
	}
}

// A delete scoped to the wrong student is a 404 and the row survives.
func TestDeleteExtracurricular_WrongStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	created := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), minimalCreateInput("Mine", "2025-01-01"))

	err := repo.DeleteExtracurricular(context.Background(), created.ID, uuid.NewString())
	assertNotFound(t, err)

	if n := countRows(t, db, "extracurriculars", created.ID); n != 1 {
		t.Errorf("row was deleted by another student's request")
	}
}

// Deleting a missing id is a 404, and deleting twice is a 404 the second time.
func TestDeleteExtracurricular_NotFoundAndIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	assertNotFound(t, repo.DeleteExtracurricular(ctx, 123456, uuid.NewString()))

	studentID := uuid.NewString()
	created := createTestExtracurricular(t, repo, studentID, createTestUser(t, db), minimalCreateInput("Twice", "2025-01-01"))
	if err := repo.DeleteExtracurricular(ctx, created.ID, studentID); err != nil {
		t.Fatalf("first delete failed: %v", err)
	}
	assertNotFound(t, repo.DeleteExtracurricular(ctx, created.ID, studentID))
}
