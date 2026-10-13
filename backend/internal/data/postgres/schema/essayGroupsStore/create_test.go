package essayGroupRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

// test normal essay group creation with valid input
func TestCreateEssayGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayGroupRepository(db)
	ctx := context.Background()

	studentID := uuid.New()
	description := "Essays for the Common App"

	input := models.CreateEssayGroupBody{
		Name:        "Common App",
		Description: &description,
		StudentID:   studentID,
	}

	created, err := repo.CreateEssayGroup(ctx, input)
	if err != nil {
		t.Fatalf("CreateEssayGroup failed: %v", err)
	}

	if created.ID == uuid.Nil {
		t.Error("expected a generated ID, got the nil UUID")
	}
	if created.Name != input.Name {
		t.Errorf("expected Name %q, got %q", input.Name, created.Name)
	}
	if created.StudentID != studentID {
		t.Errorf("expected StudentID %v, got %v", studentID, created.StudentID)
	}
	if created.Description == nil {
		t.Fatal("expected non-nil Description")
	}
	if *created.Description != description {
		t.Errorf("expected Description %q, got %q", description, *created.Description)
	}
}

// description is nullable, so omitting it must round-trip as nil rather than erroring
func TestCreateEssayGroup_NilDescription(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayGroupRepository(db)
	ctx := context.Background()

	created, err := repo.CreateEssayGroup(ctx, models.CreateEssayGroupBody{
		Name:      "Supplementals",
		StudentID: uuid.New(),
	})
	if err != nil {
		t.Fatalf("expected create with nil description to succeed, got error: %v", err)
	}

	if created.Description != nil {
		t.Errorf("expected nil Description, got %v", *created.Description)
	}
}

// the name column is UNIQUE, so a duplicate must surface as a 409 rather than a 500
func TestCreateEssayGroup_DuplicateName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayGroupRepository(db)
	ctx := context.Background()

	input := models.CreateEssayGroupBody{
		Name:      "Duplicate Group",
		StudentID: uuid.New(),
	}

	if _, err := repo.CreateEssayGroup(ctx, input); err != nil {
		t.Fatalf("first CreateEssayGroup failed: %v", err)
	}

	// A different student reusing the name still collides, since name is globally unique.
	_, err := repo.CreateEssayGroup(ctx, models.CreateEssayGroupBody{
		Name:      "Duplicate Group",
		StudentID: uuid.New(),
	})
	if err == nil {
		t.Fatal("expected an error creating a duplicate name, got nil")
	}

	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected an errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.GetStatus() != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, httpErr.GetStatus())
	}
}
