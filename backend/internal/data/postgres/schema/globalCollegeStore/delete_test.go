package globalCollegeRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

func TestDeleteGlobalCollege(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	created, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Delete Test University",
		SchoolLocation: "Delete City, DC",
	})
	if err != nil {
		t.Fatalf("setup CreateGlobalCollege failed: %v", err)
	}

	err = repo.DeleteGlobalCollege(ctx, created.ID)
	if err != nil {
		t.Fatalf("DeleteGlobalCollege failed: %v", err)
	}

	_, err = repo.GetGlobalCollege(ctx, created.ID)
	if err == nil {
		t.Error("expected an error fetching a deleted college, got nil")
	}
}

func TestDeleteGlobalCollege_NotFound(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	err := repo.DeleteGlobalCollege(ctx, 999999)
	if err == nil {
		t.Fatal("expected an error deleting a nonexistent college, got nil")
	}
}
