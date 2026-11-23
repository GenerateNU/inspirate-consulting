package studentRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// test normal student retrieval with a valid ID, checking every scanned column.
// GetStudent scans by position, so a reordered or dropped column in get_student.sql
// shows up here as a mismatched field.
func TestGetStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	organization := "Debate Club"
	interest := "Computer Science"

	created := createTestStudent(t, db, models.Student{
		UserID:         createTestUser(t, db, "Findable Student"),
		GraduationYear: 2028,
		Organization:   &organization,
		GPA:            3.85,
		ReviewBalance:  3,
		CounselorID:    &counselorID,
		Interest:       &interest,
	})

	fetched, err := repo.GetStudent(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetStudent failed: %v", err)
	}

	if fetched.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, fetched.ID)
	}
	if fetched.UserID != created.UserID {
		t.Errorf("expected UserID %s, got %s", created.UserID, fetched.UserID)
	}
	if fetched.GraduationYear != created.GraduationYear {
		t.Errorf("expected GraduationYear %d, got %d", created.GraduationYear, fetched.GraduationYear)
	}
	if fetched.GPA != created.GPA {
		t.Errorf("expected GPA %v, got %v", created.GPA, fetched.GPA)
	}
	if fetched.ReviewBalance != created.ReviewBalance {
		t.Errorf("expected ReviewBalance %d, got %d", created.ReviewBalance, fetched.ReviewBalance)
	}
	if fetched.Organization == nil || *fetched.Organization != organization {
		t.Errorf("expected Organization %q, got %v", organization, fetched.Organization)
	}
	if fetched.Interest == nil || *fetched.Interest != interest {
		t.Errorf("expected Interest %q, got %v", interest, fetched.Interest)
	}
	if fetched.CounselorID == nil || *fetched.CounselorID != counselorID {
		t.Errorf("expected CounselorID %s, got %v", counselorID, fetched.CounselorID)
	}
}

// organization and student_interest are nullable, so a student without them must scan back as nil
func TestGetStudent_NullableFields(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	created := createTestStudent(t, db, models.Student{
		UserID:         createTestUser(t, db, "Sparse Student"),
		GraduationYear: 2030,
		GPA:            3,
		ReviewBalance:  0,
		CounselorID:    &counselorID,
	})

	fetched, err := repo.GetStudent(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetStudent failed: %v", err)
	}

	if fetched.Organization != nil {
		t.Errorf("expected nil Organization, got %q", *fetched.Organization)
	}
	if fetched.Interest != nil {
		t.Errorf("expected nil Interest, got %q", *fetched.Interest)
	}
	if fetched.ReviewBalance != 0 {
		t.Errorf("expected ReviewBalance 0, got %d", fetched.ReviewBalance)
	}
}

// test that retrieval returns only the requested student when several exist
func TestGetStudent_ReturnsRequestedStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	target := createTestStudent(t, db, models.Student{
		UserID:         createTestUser(t, db, "Target Student"),
		GraduationYear: 2027,
		GPA:            4,
		ReviewBalance:  5,
		CounselorID:    &counselorID,
	})
	other := createTestStudent(t, db, models.Student{
		UserID:         createTestUser(t, db, "Other Student"),
		GraduationYear: 2029,
		GPA:            2,
		ReviewBalance:  1,
		CounselorID:    &counselorID,
	})

	fetched, err := repo.GetStudent(ctx, target.ID)
	if err != nil {
		t.Fatalf("GetStudent failed: %v", err)
	}

	if fetched.ID != target.ID {
		t.Fatalf("expected ID %s, got %s (other student is %s)", target.ID, fetched.ID, other.ID)
	}
	if fetched.UserID != target.UserID {
		t.Errorf("expected UserID %s, got %s", target.UserID, fetched.UserID)
	}
	if fetched.GraduationYear != target.GraduationYear {
		t.Errorf("expected GraduationYear %d, got %d", target.GraduationYear, fetched.GraduationYear)
	}
	if fetched.ReviewBalance != target.ReviewBalance {
		t.Errorf("expected ReviewBalance %d, got %d", target.ReviewBalance, fetched.ReviewBalance)
	}
}

// test retrieval of a student with a nonexistent ID returns a 404
func TestGetStudent_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	fetched, err := repo.GetStudent(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
	if fetched != nil {
		t.Errorf("expected nil student on error, got %+v", fetched)
	}
	if httpErr, ok := err.(errs.HTTPError); !ok || httpErr.Code != 404 {
		t.Errorf("expected a 404 HTTPError, got %v", err)
	}
}
