package essayRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

// Regression test: models.Essays has 7 fields and pgx.RowToStructByPos requires the
// query to select exactly that many columns in the same order. Dropping essay_group_id
// from the SELECT makes this fail at runtime.
func TestGetEssaysFromStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	studentID := uuid.New()
	groupID := insertEssayGroup(t, db, studentID, "Common App")

	err := repo.CreateEssay(ctx, models.Essays{
		StudentID:     studentID,
		Type:          "personal-statement",
		LinkToContent: "https://docs.google.com/document/d/abc123",
		EssayGroupID:  &groupID,
	})
	if err != nil {
		t.Fatalf("setup CreateEssay failed: %v", err)
	}

	essays, err := repo.GetEssaysFromStudent(ctx, studentID)
	if err != nil {
		t.Fatalf("GetEssaysFromStudent failed: %v", err)
	}

	if len(essays) != 1 {
		t.Fatalf("expected 1 essay, got %d", len(essays))
	}

	got := essays[0]
	if got.StudentID != studentID {
		t.Errorf("expected StudentID %v, got %v", studentID, got.StudentID)
	}
	if got.Type != "personal-statement" {
		t.Errorf("expected Type %q, got %q", "personal-statement", got.Type)
	}
	if got.Status != models.Draft {
		t.Errorf("expected the database default status %q, got %q", models.Draft, got.Status)
	}
	if got.EssayGroupID == nil {
		t.Fatal("expected EssayGroupID to be populated")
	}
	if *got.EssayGroupID != groupID {
		t.Errorf("expected EssayGroupID %v, got %v", groupID, *got.EssayGroupID)
	}
}

// essay_group_id is nullable, so an ungrouped essay must scan back as nil
func TestGetEssaysFromStudent_NoGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	studentID := uuid.New()

	err := repo.CreateEssay(ctx, models.Essays{
		StudentID:     studentID,
		Type:          "supplemental",
		LinkToContent: "https://docs.google.com/document/d/no-group",
	})
	if err != nil {
		t.Fatalf("setup CreateEssay failed: %v", err)
	}

	essays, err := repo.GetEssaysFromStudent(ctx, studentID)
	if err != nil {
		t.Fatalf("GetEssaysFromStudent failed: %v", err)
	}

	if len(essays) != 1 {
		t.Fatalf("expected 1 essay, got %d", len(essays))
	}
	if essays[0].EssayGroupID != nil {
		t.Errorf("expected nil EssayGroupID, got %v", *essays[0].EssayGroupID)
	}
}

// test that filtering by group returns only that group's essays
func TestGetEssaysByGroup(t *testing.T) {
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

	essaysToCreate := []models.Essays{
		{StudentID: studentID, Type: "personal-statement", LinkToContent: "https://example.com/1", EssayGroupID: &groupID},
		{StudentID: studentID, Type: "why-us", LinkToContent: "https://example.com/2", EssayGroupID: &groupID},
		{StudentID: studentID, Type: "supplemental", LinkToContent: "https://example.com/3", EssayGroupID: &otherGroupID},
		// Ungrouped essays must not appear in any group's result.
		{StudentID: studentID, Type: "orphan", LinkToContent: "https://example.com/4"},
	}

	for _, essay := range essaysToCreate {
		if err := repo.CreateEssay(ctx, essay); err != nil {
			t.Fatalf("setup CreateEssay(%q) failed: %v", essay.Type, err)
		}
	}

	essays, err := repo.GetEssaysByGroup(ctx, groupID)
	if err != nil {
		t.Fatalf("GetEssaysByGroup failed: %v", err)
	}

	if len(essays) != 2 {
		t.Fatalf("expected 2 essays in the group, got %d", len(essays))
	}

	for _, e := range essays {
		if e.EssayGroupID == nil || *e.EssayGroupID != groupID {
			t.Errorf("expected every essay to belong to group %v, got %v", groupID, e.EssayGroupID)
		}
	}
}

// an unknown group is not an error, it is simply an empty result
func TestGetEssaysByGroup_UnknownGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayRepository(db)
	ctx := context.Background()

	essays, err := repo.GetEssaysByGroup(ctx, uuid.New())
	if err != nil {
		t.Fatalf("GetEssaysByGroup failed: %v", err)
	}

	if len(essays) != 0 {
		t.Errorf("expected no essays, got %d", len(essays))
	}
}
