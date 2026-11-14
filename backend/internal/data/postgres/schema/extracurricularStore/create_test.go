package extracurricularRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

// Every field, including the application-detail columns, is persisted and returned.
func TestCreateExtracurricular_AllFields(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	studentID := uuid.NewString()
	userID := createTestUser(t, db)

	input := models.CreateExtracurricularInput{
		Body: models.CreateExtracurricularRequest{
			Name:                   "Varsity Debate",
			Status:                 models.ExtracurricularStatusHaveDone,
			Type:                   models.ExtracurricularTypeMaintenance,
			Description:            ptr("Competed in regional tournaments"),
			LeadershipRole:         ptr("Captain"),
			StartDate:              "2024-09-01",
			EndDate:                ptr("2025-06-15"),
			Organization:           "Lincoln High",
			HoursPerWeek:           ptr(7.5),
			WeeksPerYear:           ptr(int16(36)),
			Goals:                  ptr("Qualify for nationals"),
			ReadyForApplication:    ptr(true),
			ApplicationDescription: ptr("Led a team of 12 debaters to a state title."),
		},
	}

	out, err := repo.CreateExtracurricular(ctx, studentID, userID, input)
	if err != nil {
		t.Fatalf("CreateExtracurricular failed: %v", err)
	}
	got := out.Body

	if got.ID == 0 {
		t.Error("expected a generated ID")
	}
	if got.StudentID != studentID {
		t.Errorf("StudentID = %q, want %q", got.StudentID, studentID)
	}
	if got.UserID != userID {
		t.Errorf("UserID = %q, want %q", got.UserID, userID)
	}
	if got.Name != "Varsity Debate" {
		t.Errorf("Name = %q", got.Name)
	}
	if got.Status != models.ExtracurricularStatusHaveDone {
		t.Errorf("Status = %q", got.Status)
	}
	if got.Type != models.ExtracurricularTypeMaintenance {
		t.Errorf("Type = %q", got.Type)
	}
	if got.Description == nil || *got.Description != "Competed in regional tournaments" {
		t.Errorf("Description = %v", got.Description)
	}
	if got.LeadershipRole == nil || *got.LeadershipRole != "Captain" {
		t.Errorf("LeadershipRole = %v", got.LeadershipRole)
	}
	if !got.StartDate.Equal(time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("StartDate = %v", got.StartDate)
	}
	if got.EndDate == nil || !got.EndDate.Equal(time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("EndDate = %v", got.EndDate)
	}
	if got.Organization != "Lincoln High" {
		t.Errorf("Organization = %q", got.Organization)
	}
	if got.HoursPerWeek == nil || *got.HoursPerWeek != 7.5 {
		t.Errorf("HoursPerWeek = %v, want 7.5", got.HoursPerWeek)
	}
	if got.WeeksPerYear == nil || *got.WeeksPerYear != 36 {
		t.Errorf("WeeksPerYear = %v, want 36", got.WeeksPerYear)
	}
	if got.Goals == nil || *got.Goals != "Qualify for nationals" {
		t.Errorf("Goals = %v", got.Goals)
	}
	if !got.ReadyForApplication {
		t.Error("ReadyForApplication = false, want true")
	}
	if got.ApplicationDescription == nil || *got.ApplicationDescription != "Led a team of 12 debaters to a state title." {
		t.Errorf("ApplicationDescription = %v", got.ApplicationDescription)
	}
	if got.CreatedAt.IsZero() || got.ModifiedAt.IsZero() {
		t.Error("expected CreatedAt and ModifiedAt to be set")
	}
}

// Omitting every optional field stores NULLs and defaults ready_for_application to false.
func TestCreateExtracurricular_OptionalFieldsOmitted(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)

	got := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), minimalCreateInput("Chess Club", "2025-01-10"))

	if got.Description != nil {
		t.Errorf("Description = %v, want nil (column is nullable)", *got.Description)
	}
	if got.LeadershipRole != nil {
		t.Errorf("LeadershipRole = %v, want nil", *got.LeadershipRole)
	}
	if got.EndDate != nil {
		t.Errorf("EndDate = %v, want nil", *got.EndDate)
	}
	if got.HoursPerWeek != nil {
		t.Errorf("HoursPerWeek = %v, want nil", *got.HoursPerWeek)
	}
	if got.WeeksPerYear != nil {
		t.Errorf("WeeksPerYear = %v, want nil", *got.WeeksPerYear)
	}
	if got.Goals != nil {
		t.Errorf("Goals = %v, want nil", *got.Goals)
	}
	if got.ReadyForApplication {
		t.Error("ReadyForApplication = true, want default false")
	}
	if got.ApplicationDescription != nil {
		t.Errorf("ApplicationDescription = %v, want nil", *got.ApplicationDescription)
	}
}

// An explicit false for ready_for_application is stored as false.
func TestCreateExtracurricular_ReadyForApplicationExplicitFalse(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)

	input := minimalCreateInput("Band", "2025-01-10")
	input.Body.ReadyForApplication = ptr(false)

	got := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), input)
	if got.ReadyForApplication {
		t.Error("ReadyForApplication = true, want false")
	}
}

// weeks_per_year accepts both ends of the 0-52 range.
func TestCreateExtracurricular_WeeksPerYearBoundaries(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	userID := createTestUser(t, db)

	for _, weeks := range []int16{0, 52} {
		input := minimalCreateInput("Boundary", "2025-01-10")
		input.Body.WeeksPerYear = ptr(weeks)

		got := createTestExtracurricular(t, repo, uuid.NewString(), userID, input)
		if got.WeeksPerYear == nil || *got.WeeksPerYear != weeks {
			t.Errorf("WeeksPerYear = %v, want %d", got.WeeksPerYear, weeks)
		}
	}
}

// The database check constraint rejects weeks_per_year outside 0-52 even if
// request validation is bypassed.
func TestCreateExtracurricular_WeeksPerYearOutOfRange(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	userID := createTestUser(t, db)

	for _, weeks := range []int16{-1, 53} {
		input := minimalCreateInput("Out of range", "2025-01-10")
		input.Body.WeeksPerYear = ptr(weeks)

		if _, err := repo.CreateExtracurricular(context.Background(), uuid.NewString(), userID, input); err == nil {
			t.Errorf("expected check constraint error for weeks_per_year=%d, got nil", weeks)
		}
	}
}

// Fractional and zero hours_per_week round-trip through the numeric column.
func TestCreateExtracurricular_HoursPerWeekPrecision(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	userID := createTestUser(t, db)

	for _, hours := range []float64{0, 0.25, 12.75, 40} {
		input := minimalCreateInput("Hours", "2025-01-10")
		input.Body.HoursPerWeek = ptr(hours)

		got := createTestExtracurricular(t, repo, uuid.NewString(), userID, input)
		if got.HoursPerWeek == nil || *got.HoursPerWeek != hours {
			t.Errorf("HoursPerWeek = %v, want %v", got.HoursPerWeek, hours)
		}
	}
}

// Invalid dates are rejected before the query runs.
func TestCreateExtracurricular_InvalidDates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	userID := createTestUser(t, db)

	badStart := minimalCreateInput("Bad start", "09/01/2024")
	if _, err := repo.CreateExtracurricular(context.Background(), uuid.NewString(), userID, badStart); err == nil {
		t.Error("expected error for malformed start_date")
	}

	badEnd := minimalCreateInput("Bad end", "2024-09-01")
	badEnd.Body.EndDate = ptr("not-a-date")
	if _, err := repo.CreateExtracurricular(context.Background(), uuid.NewString(), userID, badEnd); err == nil {
		t.Error("expected error for malformed end_date")
	}
}

// user_id must reference an existing user.
func TestCreateExtracurricular_UnknownUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)

	_, err := repo.CreateExtracurricular(context.Background(), uuid.NewString(), uuid.NewString(), minimalCreateInput("Orphan", "2025-01-10"))
	if err == nil {
		t.Error("expected foreign key error for unknown user_id")
	}
}

// Invalid enum values are rejected by the table's check constraints.
func TestCreateExtracurricular_InvalidEnums(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	userID := createTestUser(t, db)

	badStatus := minimalCreateInput("Bad status", "2025-01-10")
	badStatus.Body.Status = "planning"
	if _, err := repo.CreateExtracurricular(context.Background(), uuid.NewString(), userID, badStatus); err == nil {
		t.Error("expected check constraint error for invalid status")
	}

	badType := minimalCreateInput("Bad type", "2025-01-10")
	badType.Body.Type = "hobby"
	if _, err := repo.CreateExtracurricular(context.Background(), uuid.NewString(), userID, badType); err == nil {
		t.Error("expected check constraint error for invalid type")
	}
}
