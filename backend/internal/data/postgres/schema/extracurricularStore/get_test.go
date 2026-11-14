package extracurricularRepository

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// The extracurricular is returned with all of its fields, achievements and notes.
func TestGetExtracurricular_WithAchievementsAndNotes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	studentID := uuid.NewString()
	userID := createTestUser(t, db)

	input := fullCreateInput()
	input.Body.ReadyForApplication = ptr(true)
	created := createTestExtracurricular(t, repo, studentID, userID, input)

	oldAch := insertAchievement(t, db, created.ID, "2024-03-01", "Won regionals")
	newAch := insertAchievement(t, db, created.ID, "2024-11-15", "Placed at state")
	oldNote := insertNote(t, db, created.ID, "2024-01-01T10:00:00Z", "Started the season")
	newNote := insertNote(t, db, created.ID, "2024-06-01T10:00:00Z", "Ask coach for a rec letter")

	got, err := repo.GetExtracurricular(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetExtracurricular failed: %v", err)
	}

	assertDetailParentFields(t, got, created, studentID, userID)
	assertDetailAchievements(t, got, created.ID, newAch, oldAch)
	assertDetailNotes(t, got, created.ID, newNote, oldNote)
}

func assertDetailParentFields(t *testing.T, got *models.ExtracurricularDetail, created models.Extracurricular, studentID, userID string) {
	t.Helper()
	// Parent fields
	if got.ID != created.ID || got.StudentID != studentID || got.UserID != userID {
		t.Errorf("identity mismatch: %+v", got.Extracurricular)
	}
	if got.Name != created.Name || got.Organization != created.Organization {
		t.Errorf("name/org mismatch: %q %q", got.Name, got.Organization)
	}
	if got.HoursPerWeek == nil || *got.HoursPerWeek != 3.0 {
		t.Errorf("HoursPerWeek = %v", got.HoursPerWeek)
	}
	if got.WeeksPerYear == nil || *got.WeeksPerYear != 30 {
		t.Errorf("WeeksPerYear = %v", got.WeeksPerYear)
	}
	if got.Goals == nil || *got.Goals != "Original goals" {
		t.Errorf("Goals = %v", got.Goals)
	}
	if !got.ReadyForApplication {
		t.Error("ReadyForApplication = false, want true")
	}
	if got.ApplicationDescription == nil || *got.ApplicationDescription != "Original application description" {
		t.Errorf("ApplicationDescription = %v", got.ApplicationDescription)
	}
}

func assertDetailAchievements(t *testing.T, got *models.ExtracurricularDetail, ecID, newAch, oldAch int64) {
	t.Helper()
	// Achievements: newest entry_date first
	if len(got.Achievements) != 2 {
		t.Fatalf("expected 2 achievements, got %d", len(got.Achievements))
	}
	if got.Achievements[0].ID != newAch || got.Achievements[1].ID != oldAch {
		t.Errorf("achievement order = [%d %d], want [%d %d]", got.Achievements[0].ID, got.Achievements[1].ID, newAch, oldAch)
	}
	a := got.Achievements[0]
	if a.ExtracurricularID != ecID {
		t.Errorf("achievement ExtracurricularID = %d", a.ExtracurricularID)
	}
	if !a.EntryDate.Equal(time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("achievement EntryDate = %v", a.EntryDate)
	}
	if a.Body != "Placed at state" {
		t.Errorf("achievement Body = %q", a.Body)
	}
	if a.CreatedAt.IsZero() || a.ModifiedAt.IsZero() {
		t.Error("achievement timestamps not set")
	}
}

func assertDetailNotes(t *testing.T, got *models.ExtracurricularDetail, ecID, newNote, oldNote int64) {
	t.Helper()
	// Notes: newest created_at first
	if len(got.Notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(got.Notes))
	}
	if got.Notes[0].ID != newNote || got.Notes[1].ID != oldNote {
		t.Errorf("note order = [%d %d], want [%d %d]", got.Notes[0].ID, got.Notes[1].ID, newNote, oldNote)
	}
	n := got.Notes[0]
	if n.ExtracurricularID != ecID || n.Body != "Ask coach for a rec letter" {
		t.Errorf("note = %+v", n)
	}
	if !n.CreatedAt.Equal(time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)) {
		t.Errorf("note CreatedAt = %v", n.CreatedAt)
	}
}

// With no children, both lists are empty, non-nil slices so JSON shows [].
func TestGetExtracurricular_NoChildren(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	created := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), minimalCreateInput("Lonely", "2025-01-01"))

	got, err := repo.GetExtracurricular(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetExtracurricular failed: %v", err)
	}
	if got.Achievements == nil || len(got.Achievements) != 0 {
		t.Errorf("Achievements = %#v, want empty non-nil slice", got.Achievements)
	}
	if got.Notes == nil || len(got.Notes) != 0 {
		t.Errorf("Notes = %#v, want empty non-nil slice", got.Notes)
	}
	if got.Description != nil || got.ReadyForApplication {
		t.Errorf("unexpected defaults: %+v", got.Extracurricular)
	}
}

// Children of other extracurriculars, even for the same student, are excluded.
func TestGetExtracurricular_OnlyOwnChildren(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	studentID := uuid.NewString()
	userID := createTestUser(t, db)
	target := createTestExtracurricular(t, repo, studentID, userID, minimalCreateInput("Target", "2025-01-01"))
	other := createTestExtracurricular(t, repo, studentID, userID, minimalCreateInput("Other", "2024-01-01"))

	insertAchievement(t, db, target.ID, "2025-02-01", "Mine")
	insertAchievement(t, db, other.ID, "2025-02-01", "Not mine")
	insertNote(t, db, target.ID, "2025-02-01T00:00:00Z", "Mine")
	insertNote(t, db, other.ID, "2025-02-01T00:00:00Z", "Not mine")
	insertNote(t, db, other.ID, "2025-02-02T00:00:00Z", "Also not mine")

	got, err := repo.GetExtracurricular(context.Background(), target.ID)
	if err != nil {
		t.Fatalf("GetExtracurricular failed: %v", err)
	}
	if len(got.Achievements) != 1 || got.Achievements[0].Body != "Mine" {
		t.Errorf("Achievements = %+v", got.Achievements)
	}
	if len(got.Notes) != 1 || got.Notes[0].Body != "Mine" {
		t.Errorf("Notes = %+v", got.Notes)
	}
}

// Achievements on the same date fall back to newest id first.
func TestGetExtracurricular_AchievementTieBreak(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	created := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), minimalCreateInput("Ties", "2025-01-01"))

	first := insertAchievement(t, db, created.ID, "2025-05-05", "First")
	second := insertAchievement(t, db, created.ID, "2025-05-05", "Second")

	got, err := repo.GetExtracurricular(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetExtracurricular failed: %v", err)
	}
	if len(got.Achievements) != 2 || got.Achievements[0].ID != second || got.Achievements[1].ID != first {
		t.Errorf("tie-break order wrong: %+v", got.Achievements)
	}
}

// A missing id returns a 404 HTTPError.
func TestGetExtracurricular_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)

	got, err := repo.GetExtracurricular(context.Background(), 987654)
	if got != nil {
		t.Errorf("expected nil result, got %+v", got)
	}
	var httpErr errs.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != http.StatusNotFound {
		t.Errorf("expected 404 HTTPError, got %v", err)
	}
}
