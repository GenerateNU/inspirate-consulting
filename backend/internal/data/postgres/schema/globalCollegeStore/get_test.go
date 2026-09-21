package globalCollegeRepository

import (
	"context"
	"testing"

	"inspirate-consulting/internal/models"
)

// test normal global college retrieval with valid input/ID
func TestGetGlobalCollege(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	input := models.CreateGlobalCollegeRequestBody{
		SchoolName: "Findable University",
		SchoolLocation: "Search City, SC",
	}
	createOutput, err := repo.CreateGlobalCollege(ctx, input)
	if err != nil {
		t.Fatalf("setup CreateGlobalCollege failed: %v", err)
	}
	created := createOutput
	cleanupCollege(t, db, created.ID)

	getOutput, err := repo.GetGlobalCollege(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetGlobalCollege failed: %v", err)
	}
	fetched := getOutput

	if fetched.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, fetched.ID)
	}
	if fetched.SchoolName != created.SchoolName {
		t.Errorf("expected SchoolName %q, got %q", created.SchoolName, fetched.SchoolName)
	}
}

// test retrieval of a global college with a nonexistent ID returns an error
func TestGetGlobalCollege_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetGlobalCollege(ctx, -1)
	if err == nil {
		t.Fatal("expected an error for a nonexistent id, got nil")
	}
}