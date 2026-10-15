package globalCollegeRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

func TestUpdateGlobalCollege(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	created, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Update Test University",
		SchoolLocation: "Update City, UC",
	})
	if err != nil {
		t.Fatalf("setup CreateGlobalCollege failed: %v", err)
	}

	updated, err := repo.UpdateGlobalCollege(ctx, created.ID, models.UpdateGlobalCollegeRequestBody{
		SchoolName:     "Updated Name University",
		SchoolLocation: "New City, NC",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("UpdateGlobalCollege failed: %v", err)
	}

	if updated.SchoolName != "Updated Name University" {
		t.Errorf("expected SchoolName 'Updated Name University', got %q", updated.SchoolName)
	}
	if updated.SchoolLocation != "New City, NC" {
		t.Errorf("expected SchoolLocation 'New City, NC', got %q", updated.SchoolLocation)
	}
}

func TestUpdateGlobalCollege_NotFound(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.UpdateGlobalCollege(ctx, 999999, models.UpdateGlobalCollegeRequestBody{
		SchoolName:     "Nonexistent University",
		SchoolLocation: "Nowhere, NW",
	}, nearDuplicateThreshold)
	if err == nil {
		t.Fatal("expected an error updating a nonexistent college, got nil")
	}
}
