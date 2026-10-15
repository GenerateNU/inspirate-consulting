package personalCollegeApplicationRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

func TestDeletePersonalCollegeApplication(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()

	collegeID := createTestGlobalCollege(t, db, "Delete Test University", "Delete City, DC")
	studentID := "00000000-0000-0000-0000-000000000011"

	created, err := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID,
		ApplicationType: "ED",
		Category:        "reach",
	})
	if err != nil {
		t.Fatalf("setup CreatePersonalCollegeApplication failed: %v", err)
	}

	err = repo.DeletePersonalCollegeApplication(ctx, studentID, created.ID)
	if err != nil {
		t.Fatalf("DeletePersonalCollegeApplication failed: %v", err)
	}

	results, err := repo.ListPersonalCollegeApplicationsByStudentID(ctx, studentID)
	if err != nil {
		t.Fatalf("ListPersonalCollegeApplicationsByStudentID after delete failed: %v", err)
	}
	for _, r := range results {
		if r.ID == created.ID {
			t.Error("expected deleted application to no longer appear in list")
		}
	}
}

func TestDeletePersonalCollegeApplication_NotFound(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()

	err := repo.DeletePersonalCollegeApplication(ctx, "00000000-0000-0000-0000-000000000012", 999999)
	if err == nil {
		t.Fatal("expected an error deleting a nonexistent application, got nil")
	}
}

func TestDeletePersonalCollegeApplication_WrongStudentFails(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()

	collegeID := createTestGlobalCollege(t, db, "Ownership Test University", "Ownership City, OC")
	ownerStudentID := "00000000-0000-0000-0000-000000000013"
	otherStudentID := "00000000-0000-0000-0000-000000000014"

	created, err := repo.CreatePersonalCollegeApplication(ctx, ownerStudentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID,
		ApplicationType: "ED",
		Category:        "reach",
	})
	if err != nil {
		t.Fatalf("setup CreatePersonalCollegeApplication failed: %v", err)
	}

	// Attempting to delete another student's application should fail
	err = repo.DeletePersonalCollegeApplication(ctx, otherStudentID, created.ID)
	if err == nil {
		t.Fatal("expected an error deleting another student's application, got nil")
	}

	// Confirm it's still there for the real owner.
	results, err := repo.ListPersonalCollegeApplicationsByStudentID(ctx, ownerStudentID)
	if err != nil {
		t.Fatalf("ListPersonalCollegeApplicationsByStudentID failed: %v", err)
	}
	found := false
	for _, r := range results {
		if r.ID == created.ID {
			found = true
		}
	}
	if !found {
		t.Error("expected the application to still exist for its actual owner")
	}
}
