package essayRepository

import (
	"context"
	"testing"

	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// insertEssayGroup creates a group directly so the essay store tests do not depend
// on the essay group repository package.
func insertEssayGroup(t *testing.T, db *pgxpool.Pool, studentID uuid.UUID, name string) uuid.UUID {
	t.Helper()

	var groupID uuid.UUID
	err := db.QueryRow(
		context.Background(),
		`INSERT INTO public.essay_groups (name, student_id) VALUES ($1, $2) RETURNING id`,
		name,
		studentID,
	).Scan(&groupID)
	if err != nil {
		t.Fatalf("setup insertEssayGroup(%q) failed: %v", name, err)
	}

	return groupID
}

// createEssay inserts an essay and returns it, since CreateEssay itself returns only an error.
func createEssay(t *testing.T, repo *EssayRepository, studentID uuid.UUID, body models.CreateEssayBody) models.Essays {
	t.Helper()

	ctx := context.Background()
	if err := repo.CreateEssay(ctx, studentID, body); err != nil {
		t.Fatalf("setup CreateEssay failed: %v", err)
	}

	essays, err := repo.GetEssaysFromStudent(ctx, studentID)
	if err != nil {
		t.Fatalf("setup GetEssaysFromStudent failed: %v", err)
	}
	if len(essays) != 1 {
		t.Fatalf("setup expected exactly 1 essay for the student, got %d", len(essays))
	}

	return essays[0]
}
