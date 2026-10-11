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
	}, nearDuplicateThreshold)
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

// does not reject itself as a duplicate when updating with no name/location changes
func TestUpdateGlobalCollege_NoChangeDoesNotFalsePositiveOnSelf(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	created, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Self Match University",
		SchoolLocation: "Self City, SC",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("setup CreateGlobalCollege failed: %v", err)
	}

	updated, err := repo.UpdateGlobalCollege(ctx, created.ID, models.UpdateGlobalCollegeRequestBody{
		SchoolName:     "Self Match University",
		SchoolLocation: "Self City, SC",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("expected no error updating a college with its own unchanged name, got: %v", err)
	}
	if updated.SchoolName != "Self Match University" {
		t.Errorf("expected SchoolName unchanged, got %q", updated.SchoolName)
	}
}

// disallows updating a college to a name/location that is a near-duplicate of another college, unless ConfirmDuplicate is set
func TestUpdateGlobalCollege_RejectsNearDuplicateOfAnotherCollege(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewGlobalCollegeRepository(db)
	ctx := context.Background()

	_, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Northeastern University - Boston",
		SchoolLocation: "Boston, MA",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("setup create original failed: %v", err)
	}

	toUpdate, err := repo.CreateGlobalCollege(ctx, models.CreateGlobalCollegeRequestBody{
		SchoolName:     "Totally Different Name University",
		SchoolLocation: "Different City, DC",
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("setup create second college failed: %v", err)
	}

	// renaming toUpdate to collide (near-duplicate) with the first college
	_, err = repo.UpdateGlobalCollege(ctx, toUpdate.ID, models.UpdateGlobalCollegeRequestBody{
		SchoolName:     "Northeastern Boston Campus",
		SchoolLocation: "Boston, MA",
	}, nearDuplicateThreshold)
	if err == nil {
		t.Fatal("expected an error updating into a near-duplicate name/location, got nil")
	}

	// With ConfirmDuplicate set, the same update should succeed.
	updated, err := repo.UpdateGlobalCollege(ctx, toUpdate.ID, models.UpdateGlobalCollegeRequestBody{
		SchoolName:       "Northeastern Boston Campus",
		SchoolLocation:   "Boston, MA",
		ConfirmDuplicate: true,
	}, nearDuplicateThreshold)
	if err != nil {
		t.Fatalf("expected ConfirmDuplicate to allow the update, got error: %v", err)
	}
	if updated.SchoolName != "Northeastern Boston Campus" {
		t.Errorf("expected SchoolName to be updated despite the duplicate warning, got %q", updated.SchoolName)
	}
}
