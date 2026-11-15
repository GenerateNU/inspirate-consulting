package extracurricularRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// insertUser creates a user row to satisfy the extracurriculars.user_id foreign key.
func insertUser(t *testing.T, db *pgxpool.Pool) uuid.UUID {
	t.Helper()

	var userID uuid.UUID
	err := db.QueryRow(
		context.Background(),
		`INSERT INTO public.users (name, email) VALUES ($1, $2) RETURNING id`,
		"extracurricular test",
		uuid.NewString()+"@test.com",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("setup insertUser failed: %v", err)
	}

	return userID
}

// createExtracurricular inserts an extracurricular for studentID via the repository.
func createExtracurricular(t *testing.T, repo *ExtracurricularRepository, studentID, userID uuid.UUID, name string) models.Extracurricular {
	t.Helper()

	var input models.CreateExtracurricularInput
	input.Body.Name = name
	input.Body.Status = models.ExtracurricularStatusDoing
	input.Body.Type = models.ExtracurricularTypeInvestment
	desc := "test description"
	input.Body.Description = &desc
	input.Body.StartDate = "2025-09-01"
	input.Body.Organization = "Test Org"

	out, err := repo.CreateExtracurricular(context.Background(), studentID.String(), userID.String(), input)
	if err != nil {
		t.Fatalf("setup CreateExtracurricular(%q) failed: %v", name, err)
	}

	return out.Body
}

func TestGetExtracurricularStudentID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewExtracurricularRepository(db)
	ctx := context.Background()

	userID := insertUser(t, db)
	studentID := uuid.New()
	otherStudentID := uuid.New()

	created := createExtracurricular(t, repo, studentID, userID, "Robotics")
	// A second student's activity ensures the lookup is keyed on id, not just the first row.
	other := createExtracurricular(t, repo, otherStudentID, userID, "Debate")

	got, err := repo.GetExtracurricularStudentID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetExtracurricularStudentID failed: %v", err)
	}
	if got != studentID {
		t.Errorf("expected student ID %v, got %v", studentID, got)
	}

	gotOther, err := repo.GetExtracurricularStudentID(ctx, other.ID)
	if err != nil {
		t.Fatalf("GetExtracurricularStudentID(other) failed: %v", err)
	}
	if gotOther != otherStudentID {
		t.Errorf("expected student ID %v, got %v", otherStudentID, gotOther)
	}
}

func TestGetExtracurricularStudentID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewExtracurricularRepository(db)

	got, err := repo.GetExtracurricularStudentID(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected an error for a nonexistent extracurricular, got nil")
	}
	if got != uuid.Nil {
		t.Errorf("expected uuid.Nil on error, got %v", got)
	}

	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected an errs.HTTPError, got %T: %v", err, err)
	}
	if httpErr.GetStatus() != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, httpErr.GetStatus())
	}
}
