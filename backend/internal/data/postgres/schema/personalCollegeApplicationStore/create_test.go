package personalCollegeApplicationRepository

import (
	"context"
	"testing"

	"inspirate-consulting/internal/models"
)

const testStudentID = "00000000-0000-0000-0000-000000000002"

// test normal application creation with valid input
func TestCreatePersonalCollegeApplication(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	collegeID := createTestGlobalCollege(t, db, "Create Test University", "Create City, CT")

	input := models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID,
		ApplicationType: "ED",
		Category: "reach",
	}

	created, err := repo.CreatePersonalCollegeApplication(ctx, testStudentID, input)
	if err != nil {
		t.Fatalf("CreatePersonalCollegeApplication failed: %v", err)
	}
	cleanupApplication(t, db, created.ID)

	if created.ID == 0 {
		t.Error("expected a generated ID, got 0")
	}
	if created.StudentID != testStudentID {
		t.Errorf("expected StudentID %q, got %q", testStudentID, created.StudentID)
	}
	if created.GlobalCollegeID != collegeID {
		t.Errorf("expected GlobalCollegeID %d, got %d", collegeID, created.GlobalCollegeID)
	}
	if created.ApplicationType != "ED" {
		t.Errorf("expected ApplicationType %q, got %q", "ED", created.ApplicationType)
	}
	if created.Category != "reach" {
		t.Errorf("expected Category %q, got %q", "reach", created.Category)
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if created.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

// to test application creation with an invalid application_type
func TestCreatePersonalCollegeApplication_InvalidApplicationType(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	collegeID := createTestGlobalCollege(t, db, "Invalid Type University", "Invalid City, IV")

	input := models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID,
		ApplicationType: "NOT_A_REAL_TYPE",
		Category: "reach",
	}

	_, err := repo.CreatePersonalCollegeApplication(ctx, testStudentID, input)
	if err == nil {
		t.Fatal("expected an error for an invalid application_type, got nil")
	}
}

// to test application creation with an invalid category
func TestCreatePersonalCollegeApplication_InvalidCategory(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	collegeID := createTestGlobalCollege(t, db, "Invalid Category University", "Invalid City, IC")

	input := models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID,
		ApplicationType: "ED",
		Category: "not_a_real_category",
	}

	_, err := repo.CreatePersonalCollegeApplication(ctx, testStudentID, input)
	if err == nil {
		t.Fatal("expected an error for an invalid category, got nil")
	}
}

func TestCreatePersonalCollegeApplication_NonexistentGlobalCollege(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	input := models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: -999999,
		ApplicationType: "ED",
		Category: "reach",
	}

	_, err := repo.CreatePersonalCollegeApplication(ctx, testStudentID, input)
	if err == nil {
		t.Fatal("expected an error for a nonexistent global_college_id, got nil")
	}
}