package essayGroupRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

// test that listing is scoped to the requested student and does not leak other students' groups
func TestListEssayGroups(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayGroupRepository(db)
	ctx := context.Background()

	studentID := uuid.New()
	otherStudentID := uuid.New()

	inputs := []models.CreateEssayGroupBody{
		{Name: "Common App", StudentID: studentID},
		{Name: "Supplementals", StudentID: studentID},
		{Name: "Scholarships", StudentID: otherStudentID},
	}

	for _, input := range inputs {
		if _, err := repo.CreateEssayGroup(ctx, input); err != nil {
			t.Fatalf("setup CreateEssayGroup(%q) failed: %v", input.Name, err)
		}
	}

	groups, err := repo.ListEssayGroups(ctx, studentID)
	if err != nil {
		t.Fatalf("ListEssayGroups failed: %v", err)
	}

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups for the student, got %d", len(groups))
	}

	names := make(map[string]bool, len(groups))
	for _, g := range groups {
		if g.StudentID != studentID {
			t.Errorf("expected only groups for student %v, got one for %v", studentID, g.StudentID)
		}
		names[g.Name] = true
	}

	for _, want := range []string{"Common App", "Supplementals"} {
		if !names[want] {
			t.Errorf("expected group %q in the result", want)
		}
	}
	if names["Scholarships"] {
		t.Error("expected the other student's group to be excluded")
	}
}

// a student with no groups gets an empty list, not an error
func TestListEssayGroups_NoGroups(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewEssayGroupRepository(db)
	ctx := context.Background()

	groups, err := repo.ListEssayGroups(ctx, uuid.New())
	if err != nil {
		t.Fatalf("ListEssayGroups failed: %v", err)
	}

	if len(groups) != 0 {
		t.Errorf("expected no groups, got %d", len(groups))
	}
}
