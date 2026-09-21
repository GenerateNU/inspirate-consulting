package personalCollegeApplicationRepository

import (
	"context"
	"testing"

	"inspirate-consulting/internal/models"
)

// test list function returns all student college applications
// intentionally is not run in parallel, because it depends on the state of the database and could be affected by other tests
func TestListPersonalCollegeApplicationsByStudentID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := "00000000-0000-0000-0000-000000000099"
	otherStudentID := "00000000-0000-0000-0000-000000000098"

	collegeA := createTestGlobalCollege(t, db, "List Test Alpha University", "Alpha City, AA")
	collegeB := createTestGlobalCollege(t, db, "List Test Beta University", "Beta City, BB")

	appA, err := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeA,
		ApplicationType: "EA",
		Category: "target",
	})
	if err != nil {
		t.Fatalf("setup CreatePersonalCollegeApplication failed: %v", err)
	}
	cleanupApplication(t, db, appA.ID)

	appB, err := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeB,
		ApplicationType: "RD",
		Category: "safety",
	})
	if err != nil {
		t.Fatalf("setup CreatePersonalCollegeApplication failed: %v", err)
	}
	cleanupApplication(t, db, appB.ID)

	// A different student's application should not show up in studentID's results.
	otherApp, err := repo.CreatePersonalCollegeApplication(ctx, otherStudentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeA,
		ApplicationType: "ED",
		Category: "reach",
	})
	if err != nil {
		t.Fatalf("setup CreatePersonalCollegeApplication (other student) failed: %v", err)
	}
	cleanupApplication(t, db, otherApp.ID)

	results, err := repo.ListPersonalCollegeApplicationsByStudentID(ctx, studentID)
	if err != nil {
		t.Fatalf("ListPersonalCollegeApplicationsByStudentID failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 applications for studentID, got %d", len(results))
	}

	foundIDs := map[int64]bool{}
	for _, r := range results {
		foundIDs[r.ID] = true
		if r.StudentID != studentID {
			t.Errorf("expected all results to have StudentID %q, got %q", studentID, r.StudentID)
		}
	}
	if !foundIDs[appA.ID] || !foundIDs[appB.ID] {
		t.Error("expected both of studentID's applications to be present in results")
	}
	if foundIDs[otherApp.ID] {
		t.Error("expected the other student's application to NOT be present in results")
	}
}

func TestListPersonalCollegeApplicationsByStudentID_NoApplications(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	results, err := repo.ListPersonalCollegeApplicationsByStudentID(ctx, "00000000-0000-0000-0000-000000000077")
	if err != nil {
		t.Fatalf("ListPersonalCollegeApplicationsByStudentID failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 applications for a student with none, got %d", len(results))
	}
}
