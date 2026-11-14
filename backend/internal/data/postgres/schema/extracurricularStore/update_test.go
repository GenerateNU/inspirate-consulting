package extracurricularRepository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/models"
)

// fullCreateInput returns a create payload with every field populated.
func fullCreateInput() models.CreateExtracurricularInput {
	return models.CreateExtracurricularInput{
		Body: models.CreateExtracurricularRequest{
			Name:                   "Original name",
			Status:                 models.ExtracurricularStatusDoing,
			Type:                   models.ExtracurricularTypeInvestment,
			Description:            ptr("Original description"),
			LeadershipRole:         ptr("Member"),
			StartDate:              "2024-01-01",
			EndDate:                ptr("2024-12-31"),
			Organization:           "Original org",
			HoursPerWeek:           ptr(3.0),
			WeeksPerYear:           ptr(int16(30)),
			Goals:                  ptr("Original goals"),
			ReadyForApplication:    ptr(false),
			ApplicationDescription: ptr("Original application description"),
		},
	}
}

// Every field, including the new ones, can be changed in a single update.
func TestUpdateExtracurricular_AllFields(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	studentID := uuid.NewString()
	userID := createTestUser(t, db)
	existing := createTestExtracurricular(t, repo, studentID, userID, fullCreateInput())

	update := models.UpdateExtracurricularInput{
		ID: existing.ID,
		Body: models.UpdateExtracurricularRequest{
			Name:                   ptr("New name"),
			Status:                 ptr(models.ExtracurricularStatusHaveDone),
			Type:                   ptr(models.ExtracurricularTypeMaintenance),
			Description:            models.NewNullable("New description"),
			LeadershipRole:         models.NewNullable("President"),
			StartDate:              ptr("2023-02-02"),
			EndDate:                models.NewNullable("2025-05-05"),
			Organization:           ptr("New org"),
			HoursPerWeek:           models.NewNullable(10.5),
			WeeksPerYear:           models.NewNullable(int16(52)),
			Goals:                  models.NewNullable("New goals"),
			ReadyForApplication:    ptr(true),
			ApplicationDescription: models.NewNullable("New application description"),
		},
	}

	got, err := repo.UpdateExtracurricular(context.Background(), existing.ID, update)
	if err != nil {
		t.Fatalf("UpdateExtracurricular failed: %v", err)
	}

	if got.ID != existing.ID || got.StudentID != studentID || got.UserID != userID {
		t.Errorf("identity fields changed: %+v", got)
	}
	if got.Name != "New name" {
		t.Errorf("Name = %q", got.Name)
	}
	if got.Status != models.ExtracurricularStatusHaveDone {
		t.Errorf("Status = %q", got.Status)
	}
	if got.Type != models.ExtracurricularTypeMaintenance {
		t.Errorf("Type = %q", got.Type)
	}
	if got.Description == nil || *got.Description != "New description" {
		t.Errorf("Description = %v", got.Description)
	}
	if got.LeadershipRole == nil || *got.LeadershipRole != "President" {
		t.Errorf("LeadershipRole = %v", got.LeadershipRole)
	}
	if !got.StartDate.Equal(time.Date(2023, 2, 2, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("StartDate = %v", got.StartDate)
	}
	if got.EndDate == nil || !got.EndDate.Equal(time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("EndDate = %v", got.EndDate)
	}
	if got.Organization != "New org" {
		t.Errorf("Organization = %q", got.Organization)
	}
	if got.HoursPerWeek == nil || *got.HoursPerWeek != 10.5 {
		t.Errorf("HoursPerWeek = %v", got.HoursPerWeek)
	}
	if got.WeeksPerYear == nil || *got.WeeksPerYear != 52 {
		t.Errorf("WeeksPerYear = %v", got.WeeksPerYear)
	}
	if got.Goals == nil || *got.Goals != "New goals" {
		t.Errorf("Goals = %v", got.Goals)
	}
	if !got.ReadyForApplication {
		t.Error("ReadyForApplication = false, want true")
	}
	if got.ApplicationDescription == nil || *got.ApplicationDescription != "New application description" {
		t.Errorf("ApplicationDescription = %v", got.ApplicationDescription)
	}
	if got.ModifiedAt.Before(existing.ModifiedAt) {
		t.Errorf("ModifiedAt went backwards: %v -> %v", existing.ModifiedAt, got.ModifiedAt)
	}
	if !got.CreatedAt.Equal(existing.CreatedAt) {
		t.Errorf("CreatedAt changed: %v -> %v", existing.CreatedAt, got.CreatedAt)
	}
}

// An empty update leaves every column as it was.
func TestUpdateExtracurricular_EmptyBodyKeepsValues(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	existing := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), fullCreateInput())

	got, err := repo.UpdateExtracurricular(context.Background(), existing.ID, models.UpdateExtracurricularInput{ID: existing.ID})
	if err != nil {
		t.Fatalf("UpdateExtracurricular failed: %v", err)
	}

	// Only modified_at is expected to differ.
	got.ModifiedAt = existing.ModifiedAt
	if got.Name != existing.Name || got.Status != existing.Status || got.Type != existing.Type ||
		*got.Description != *existing.Description || *got.LeadershipRole != *existing.LeadershipRole ||
		!got.StartDate.Equal(existing.StartDate) || !got.EndDate.Equal(*existing.EndDate) ||
		got.Organization != existing.Organization ||
		*got.HoursPerWeek != *existing.HoursPerWeek || *got.WeeksPerYear != *existing.WeeksPerYear ||
		*got.Goals != *existing.Goals || got.ReadyForApplication != existing.ReadyForApplication ||
		*got.ApplicationDescription != *existing.ApplicationDescription {
		t.Errorf("empty update changed the row:\nbefore %+v\nafter  %+v", existing, *got)
	}
}

// Updating only the new fields leaves the original fields alone, and vice versa.
func TestUpdateExtracurricular_PartialUpdates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	existing := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), fullCreateInput())
	ctx := context.Background()

	t.Run("only application details", func(t *testing.T) {
		got, err := repo.UpdateExtracurricular(ctx, existing.ID, models.UpdateExtracurricularInput{
			ID: existing.ID,
			Body: models.UpdateExtracurricularRequest{
				HoursPerWeek: models.NewNullable(6.25),
				Goals:        models.NewNullable("Updated goals"),
			},
		})
		if err != nil {
			t.Fatalf("UpdateExtracurricular failed: %v", err)
		}
		if *got.HoursPerWeek != 6.25 || *got.Goals != "Updated goals" {
			t.Errorf("new values not applied: hours=%v goals=%v", *got.HoursPerWeek, *got.Goals)
		}
		if got.Name != existing.Name || got.Organization != existing.Organization {
			t.Errorf("unrelated fields changed: name=%q org=%q", got.Name, got.Organization)
		}
		if *got.WeeksPerYear != 30 || *got.ApplicationDescription != "Original application description" {
			t.Errorf("untouched new fields changed: weeks=%v appdesc=%v", *got.WeeksPerYear, *got.ApplicationDescription)
		}
	})

	t.Run("only original fields", func(t *testing.T) {
		got, err := repo.UpdateExtracurricular(ctx, existing.ID, models.UpdateExtracurricularInput{
			ID:   existing.ID,
			Body: models.UpdateExtracurricularRequest{Name: ptr("Renamed")},
		})
		if err != nil {
			t.Fatalf("UpdateExtracurricular failed: %v", err)
		}
		if got.Name != "Renamed" {
			t.Errorf("Name = %q", got.Name)
		}
		// Values from the previous subtest must survive.
		if *got.HoursPerWeek != 6.25 || *got.Goals != "Updated goals" {
			t.Errorf("earlier update lost: hours=%v goals=%v", *got.HoursPerWeek, *got.Goals)
		}
	})
}

// ready_for_application can be flipped both ways, including back to false.
func TestUpdateExtracurricular_ToggleReadyForApplication(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	existing := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), minimalCreateInput("Toggle", "2025-01-01"))
	ctx := context.Background()

	for _, want := range []bool{true, false, true} {
		got, err := repo.UpdateExtracurricular(ctx, existing.ID, models.UpdateExtracurricularInput{
			ID:   existing.ID,
			Body: models.UpdateExtracurricularRequest{ReadyForApplication: ptr(want)},
		})
		if err != nil {
			t.Fatalf("UpdateExtracurricular failed: %v", err)
		}
		if got.ReadyForApplication != want {
			t.Errorf("ReadyForApplication = %v, want %v", got.ReadyForApplication, want)
		}
	}
}

// New fields can be set on a row that was created without them.
func TestUpdateExtracurricular_SetFieldsThatWereNull(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	existing := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), minimalCreateInput("Bare", "2025-01-01"))

	got, err := repo.UpdateExtracurricular(context.Background(), existing.ID, models.UpdateExtracurricularInput{
		ID: existing.ID,
		Body: models.UpdateExtracurricularRequest{
			Description:            models.NewNullable("Now described"),
			WeeksPerYear:           models.NewNullable(int16(0)),
			HoursPerWeek:           models.NewNullable(0.0),
			ApplicationDescription: models.NewNullable("Now has an app blurb"),
		},
	})
	if err != nil {
		t.Fatalf("UpdateExtracurricular failed: %v", err)
	}
	if got.Description == nil || *got.Description != "Now described" {
		t.Errorf("Description = %v", got.Description)
	}
	if got.WeeksPerYear == nil || *got.WeeksPerYear != 0 {
		t.Errorf("WeeksPerYear = %v, want 0", got.WeeksPerYear)
	}
	if got.HoursPerWeek == nil || *got.HoursPerWeek != 0 {
		t.Errorf("HoursPerWeek = %v, want 0", got.HoursPerWeek)
	}
	if got.ApplicationDescription == nil || *got.ApplicationDescription != "Now has an app blurb" {
		t.Errorf("ApplicationDescription = %v", got.ApplicationDescription)
	}
	if got.Goals != nil {
		t.Errorf("Goals = %v, want still nil", *got.Goals)
	}
}

// Only the targeted row is changed.
func TestUpdateExtracurricular_DoesNotTouchOtherRows(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	studentID := uuid.NewString()
	userID := createTestUser(t, db)
	target := createTestExtracurricular(t, repo, studentID, userID, minimalCreateInput("Target", "2025-01-01"))
	bystander := createTestExtracurricular(t, repo, studentID, userID, minimalCreateInput("Bystander", "2024-01-01"))

	if _, err := repo.UpdateExtracurricular(context.Background(), target.ID, models.UpdateExtracurricularInput{
		ID:   target.ID,
		Body: models.UpdateExtracurricularRequest{Goals: models.NewNullable("Only for target"), ReadyForApplication: ptr(true)},
	}); err != nil {
		t.Fatalf("UpdateExtracurricular failed: %v", err)
	}

	list, err := repo.ListExtracurriculars(context.Background(), studentID)
	if err != nil {
		t.Fatalf("ListExtracurriculars failed: %v", err)
	}
	for _, e := range list {
		if e.ID == bystander.ID && (e.Goals != nil || e.ReadyForApplication) {
			t.Errorf("bystander row was modified: %+v", e)
		}
	}
}

// The database rejects out-of-range weeks_per_year on update too.
func TestUpdateExtracurricular_WeeksPerYearOutOfRange(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	existing := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), minimalCreateInput("Range", "2025-01-01"))

	_, err := repo.UpdateExtracurricular(context.Background(), existing.ID, models.UpdateExtracurricularInput{
		ID:   existing.ID,
		Body: models.UpdateExtracurricularRequest{WeeksPerYear: models.NewNullable(int16(60))},
	})
	if err == nil {
		t.Error("expected check constraint error for weeks_per_year=60")
	}
}

// Malformed dates are rejected before the query runs.
func TestUpdateExtracurricular_InvalidDates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	existing := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), minimalCreateInput("Dates", "2025-01-01"))
	ctx := context.Background()

	if _, err := repo.UpdateExtracurricular(ctx, existing.ID, models.UpdateExtracurricularInput{
		Body: models.UpdateExtracurricularRequest{StartDate: ptr("2025/01/01")},
	}); err == nil {
		t.Error("expected error for malformed start_date")
	}
	if _, err := repo.UpdateExtracurricular(ctx, existing.ID, models.UpdateExtracurricularInput{
		Body: models.UpdateExtracurricularRequest{EndDate: models.NewNullable("tomorrow")},
	}); err == nil {
		t.Error("expected error for malformed end_date")
	}
}

// Updating an id that does not exist returns pgx.ErrNoRows.
func TestUpdateExtracurricular_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)

	_, err := repo.UpdateExtracurricular(context.Background(), 999999, models.UpdateExtracurricularInput{
		ID:   999999,
		Body: models.UpdateExtracurricularRequest{Name: ptr("Ghost")},
	})
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("expected pgx.ErrNoRows, got %v", err)
	}
}

// An explicit null clears each nullable column; required columns and
// omitted nullable columns are left alone.
func TestUpdateExtracurricular_NullClearsNullableFields(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	existing := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), fullCreateInput())

	got, err := repo.UpdateExtracurricular(context.Background(), existing.ID, models.UpdateExtracurricularInput{
		ID: existing.ID,
		Body: models.UpdateExtracurricularRequest{
			Description:            models.Null[string](),
			LeadershipRole:         models.Null[string](),
			EndDate:                models.Null[string](),
			HoursPerWeek:           models.Null[float64](),
			WeeksPerYear:           models.Null[int16](),
			Goals:                  models.Null[string](),
			ApplicationDescription: models.Null[string](),
		},
	})
	if err != nil {
		t.Fatalf("UpdateExtracurricular failed: %v", err)
	}

	if got.Description != nil || got.LeadershipRole != nil || got.EndDate != nil ||
		got.HoursPerWeek != nil || got.WeeksPerYear != nil || got.Goals != nil ||
		got.ApplicationDescription != nil {
		t.Errorf("nullable fields not cleared: %+v", *got)
	}
	if got.Name != existing.Name || got.Organization != existing.Organization ||
		!got.StartDate.Equal(existing.StartDate) || got.ReadyForApplication != existing.ReadyForApplication {
		t.Errorf("required fields changed:\nbefore %+v\nafter  %+v", existing, *got)
	}
}

// Clearing one field leaves the other nullable fields untouched.
func TestUpdateExtracurricular_ClearOnlyEndDate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	existing := createTestExtracurricular(t, repo, uuid.NewString(), createTestUser(t, db), fullCreateInput())

	got, err := repo.UpdateExtracurricular(context.Background(), existing.ID, models.UpdateExtracurricularInput{
		ID:   existing.ID,
		Body: models.UpdateExtracurricularRequest{EndDate: models.Null[string]()},
	})
	if err != nil {
		t.Fatalf("UpdateExtracurricular failed: %v", err)
	}
	if got.EndDate != nil {
		t.Errorf("EndDate = %v, want nil", *got.EndDate)
	}
	if got.Goals == nil || *got.Goals != *existing.Goals ||
		got.HoursPerWeek == nil || *got.HoursPerWeek != *existing.HoursPerWeek ||
		got.Description == nil || *got.Description != *existing.Description {
		t.Errorf("other nullable fields changed:\nbefore %+v\nafter  %+v", existing, *got)
	}
}
