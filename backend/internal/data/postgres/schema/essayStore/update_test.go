package essayRepository

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

// Regression test: UpdateStatus returns the updated essay via RowToStructByPos, so the
// RETURNING clause must list all 7 columns of models.Essays in order.
func TestUpdateStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	studentID := uuid.New()
	groupID := insertEssayGroup(t, db, studentID, "Common App")

	created := createEssay(t, repo, models.Essays{
		StudentID:     studentID,
		Type:          "personal-statement",
		LinkToContent: "https://docs.google.com/document/d/abc123",
		EssayGroupID:  &groupID,
	})

	updated, err := repo.UpdateStatus(ctx, created.ID, models.Submitted)
	if err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}

	if updated.Status != models.Submitted {
		t.Errorf("expected Status %q, got %q", models.Submitted, updated.Status)
	}
	if updated.ID != created.ID {
		t.Errorf("expected ID %v, got %v", created.ID, updated.ID)
	}
	// Updating the status must not disturb the group membership.
	if updated.EssayGroupID == nil || *updated.EssayGroupID != groupID {
		t.Errorf("expected EssayGroupID %v to be preserved, got %v", groupID, updated.EssayGroupID)
	}
}

func TestUpdateStatus_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	_, err := repo.UpdateStatus(ctx, uuid.New(), models.Submitted)
	if err == nil {
		t.Fatal("expected an error updating a nonexistent essay, got nil")
	}

	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected an errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.GetStatus() != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, httpErr.GetStatus())
	}
}

// test adding to a group, moving between groups, and removing from a group
func TestUpdateEssayGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	studentID := uuid.New()
	groupID := insertEssayGroup(t, db, studentID, "Common App")
	otherGroupID := insertEssayGroup(t, db, studentID, "Supplementals")

	// Starts with no group.
	created := createEssay(t, repo, models.Essays{
		StudentID:     studentID,
		Type:          "personal-statement",
		LinkToContent: "https://docs.google.com/document/d/abc123",
	})
	if created.EssayGroupID != nil {
		t.Fatalf("setup expected an ungrouped essay, got group %v", *created.EssayGroupID)
	}

	// Add to a group.
	updated, err := repo.UpdateEssayGroup(ctx, created.ID, &groupID)
	if err != nil {
		t.Fatalf("UpdateEssayGroup (add) failed: %v", err)
	}
	if updated.EssayGroupID == nil || *updated.EssayGroupID != groupID {
		t.Fatalf("expected EssayGroupID %v, got %v", groupID, updated.EssayGroupID)
	}

	// Move to another group.
	updated, err = repo.UpdateEssayGroup(ctx, created.ID, &otherGroupID)
	if err != nil {
		t.Fatalf("UpdateEssayGroup (move) failed: %v", err)
	}
	if updated.EssayGroupID == nil || *updated.EssayGroupID != otherGroupID {
		t.Fatalf("expected EssayGroupID %v, got %v", otherGroupID, updated.EssayGroupID)
	}

	// Remove from its group.
	updated, err = repo.UpdateEssayGroup(ctx, created.ID, nil)
	if err != nil {
		t.Fatalf("UpdateEssayGroup (remove) failed: %v", err)
	}
	if updated.EssayGroupID != nil {
		t.Errorf("expected nil EssayGroupID, got %v", *updated.EssayGroupID)
	}

	// The rest of the essay must survive the group changes.
	if updated.Type != "personal-statement" {
		t.Errorf("expected Type to be preserved, got %q", updated.Type)
	}
	if updated.Status != models.Draft {
		t.Errorf("expected Status to be preserved, got %q", updated.Status)
	}
}

func TestUpdateEssayGroup_EssayNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	studentID := uuid.New()
	groupID := insertEssayGroup(t, db, studentID, "Common App")

	_, err := repo.UpdateEssayGroup(ctx, uuid.New(), &groupID)
	if err == nil {
		t.Fatal("expected an error updating a nonexistent essay, got nil")
	}

	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected an errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.GetStatus() != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, httpErr.GetStatus())
	}
}

// essay_group_id references essay_groups(id), so an unknown group must be a 404 not a 500
func TestUpdateEssayGroup_GroupNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	studentID := uuid.New()
	created := createEssay(t, repo, models.Essays{
		StudentID:     studentID,
		Type:          "personal-statement",
		LinkToContent: "https://docs.google.com/document/d/abc123",
	})

	unknownGroupID := uuid.New()
	_, err := repo.UpdateEssayGroup(ctx, created.ID, &unknownGroupID)
	if err == nil {
		t.Fatal("expected an error moving an essay into a nonexistent group, got nil")
	}

	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected an errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.GetStatus() != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, httpErr.GetStatus())
	}
}
