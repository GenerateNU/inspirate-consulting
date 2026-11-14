package extracurricularRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

// A student with no extracurriculars gets an empty, non-nil slice.
func TestListExtracurriculars_Empty(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)

	list, err := repo.ListExtracurriculars(context.Background(), uuid.NewString())
	if err != nil {
		t.Fatalf("ListExtracurriculars failed: %v", err)
	}
	if list == nil {
		t.Error("expected non-nil empty slice so the API returns [] not null")
	}
	if len(list) != 0 {
		t.Errorf("expected 0 items, got %d", len(list))
	}
}

// Only the requested student's rows are returned, newest start_date first.
func TestListExtracurriculars_ScopedAndOrdered(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	userID := createTestUser(t, db)
	studentID := uuid.NewString()
	otherStudentID := uuid.NewString()

	createTestExtracurricular(t, repo, studentID, userID, minimalCreateInput("Oldest", "2022-01-01"))
	createTestExtracurricular(t, repo, studentID, userID, minimalCreateInput("Newest", "2025-01-01"))
	createTestExtracurricular(t, repo, studentID, userID, minimalCreateInput("Middle", "2023-06-01"))
	createTestExtracurricular(t, repo, otherStudentID, userID, minimalCreateInput("Someone else", "2024-01-01"))

	list, err := repo.ListExtracurriculars(context.Background(), studentID)
	if err != nil {
		t.Fatalf("ListExtracurriculars failed: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 items, got %d", len(list))
	}

	want := []string{"Newest", "Middle", "Oldest"}
	for i, name := range want {
		if list[i].Name != name {
			t.Errorf("list[%d].Name = %q, want %q", i, list[i].Name, name)
		}
		if list[i].StudentID != studentID {
			t.Errorf("list[%d] belongs to student %q, want %q", i, list[i].StudentID, studentID)
		}
	}
}

// The new application-detail columns are returned by list, both set and unset.
func TestListExtracurriculars_ReturnsApplicationDetails(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	userID := createTestUser(t, db)
	studentID := uuid.NewString()

	full := minimalCreateInput("Full", "2025-01-01")
	full.Body.Description = ptr("All the details")
	full.Body.HoursPerWeek = ptr(4.5)
	full.Body.WeeksPerYear = ptr(int16(40))
	full.Body.Goals = ptr("Win state")
	full.Body.ReadyForApplication = ptr(true)
	full.Body.ApplicationDescription = ptr("Short blurb")
	createTestExtracurricular(t, repo, studentID, userID, full)
	createTestExtracurricular(t, repo, studentID, userID, minimalCreateInput("Bare", "2024-01-01"))

	list, err := repo.ListExtracurriculars(context.Background(), studentID)
	if err != nil {
		t.Fatalf("ListExtracurriculars failed: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 items, got %d", len(list))
	}

	got := list[0]
	if got.Name != "Full" {
		t.Fatalf("expected Full first, got %q", got.Name)
	}
	if got.Description == nil || *got.Description != "All the details" {
		t.Errorf("Description = %v", got.Description)
	}
	if got.HoursPerWeek == nil || *got.HoursPerWeek != 4.5 {
		t.Errorf("HoursPerWeek = %v", got.HoursPerWeek)
	}
	if got.WeeksPerYear == nil || *got.WeeksPerYear != 40 {
		t.Errorf("WeeksPerYear = %v", got.WeeksPerYear)
	}
	if got.Goals == nil || *got.Goals != "Win state" {
		t.Errorf("Goals = %v", got.Goals)
	}
	if !got.ReadyForApplication {
		t.Error("ReadyForApplication = false, want true")
	}
	if got.ApplicationDescription == nil || *got.ApplicationDescription != "Short blurb" {
		t.Errorf("ApplicationDescription = %v", got.ApplicationDescription)
	}

	bare := list[1]
	if bare.Description != nil || bare.HoursPerWeek != nil || bare.WeeksPerYear != nil ||
		bare.Goals != nil || bare.ApplicationDescription != nil || bare.ReadyForApplication {
		t.Errorf("expected bare row to have null/false application details, got %+v", bare)
	}
}

// Rows with a NULL description (allowed since the latest migration) still scan.
func TestListExtracurriculars_NullDescriptionInsertedDirectly(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	userID := createTestUser(t, db)
	studentID := uuid.NewString()

	if _, err := db.Exec(context.Background(),
		`INSERT INTO public.extracurriculars (student_id, user_id, name, status, type, start_date, organization)
		 VALUES ($1, $2, 'Direct insert', 'doing', 'investment', '2025-03-01', 'Org')`,
		studentID, userID); err != nil {
		t.Fatalf("direct insert failed: %v", err)
	}

	list, err := repo.ListExtracurriculars(context.Background(), studentID)
	if err != nil {
		t.Fatalf("ListExtracurriculars failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list))
	}
	if list[0].Description != nil {
		t.Errorf("Description = %v, want nil", *list[0].Description)
	}
	if list[0].Status != models.ExtracurricularStatusDoing {
		t.Errorf("Status = %q", list[0].Status)
	}
}
