package personalCollegeApplicationRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

func TestUpdatePersonalCollegeApplication(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()

	collegeA := createTestGlobalCollege(t, db, "Update Test Alpha University", "Alpha City, AA")
	collegeB := createTestGlobalCollege(t, db, "Update Test Beta University", "Beta City, BB")
	studentID := "00000000-0000-0000-0000-000000000021"

	created, err := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeA,
		ApplicationType: "ED",
		Category:        "reach",
	})
	if err != nil {
		t.Fatalf("setup CreatePersonalCollegeApplication failed: %v", err)
	}

	update := models.UpdatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeB,
		ApplicationType: "RD",
		Category:        "safety",
	}

	updated, err := repo.UpdatePersonalCollegeApplication(ctx, studentID, created.ID, update)
	if err != nil {
		t.Fatalf("UpdatePersonalCollegeApplication failed: %v", err)
	}

	if updated.GlobalCollegeID != collegeB {
		t.Errorf("expected GlobalCollegeID %d, got %d", collegeB, updated.GlobalCollegeID)
	}
	if updated.ApplicationType != "RD" {
		t.Errorf("expected ApplicationType RD, got %s", updated.ApplicationType)
	}
	if updated.Category != "safety" {
		t.Errorf("expected Category safety, got %s", updated.Category)
	}

	// confirm the change persisted
	fetched, err := repo.ListPersonalCollegeApplicationsByStudentID(ctx, studentID)
	if err != nil {
		t.Fatalf("ListPersonalCollegeApplicationsByStudentID failed: %v", err)
	}
	if len(fetched) != 1 || fetched[0].GlobalCollegeID != collegeB {
		t.Errorf("expected persisted GlobalCollegeID %d, got %+v", collegeB, fetched)
	}
}

func TestUpdatePersonalCollegeApplication_WrongStudentFails(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()

	collegeID := createTestGlobalCollege(t, db, "Update Ownership University", "Ownership City, OC")
	ownerStudentID := "00000000-0000-0000-0000-000000000022"
	otherStudentID := "00000000-0000-0000-0000-000000000023"

	created, err := repo.CreatePersonalCollegeApplication(ctx, ownerStudentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID,
		ApplicationType: "ED",
		Category:        "reach",
	})
	if err != nil {
		t.Fatalf("setup CreatePersonalCollegeApplication failed: %v", err)
	}

	_, err = repo.UpdatePersonalCollegeApplication(ctx, otherStudentID, created.ID, models.UpdatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID,
		ApplicationType: "RD",
		Category:        "safety",
	})
	if err == nil {
		t.Fatal("expected an error updating another student's application, got nil")
	}
}
