package studentRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// updating every field changes every field
func TestUpdateStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	created := createTestStudent(t, db, models.Student{
		UserID:        createTestUser(t, db, "Updatable Student"),
		Year:          models.Freshman,
		GPA:           3,
		ReviewBalance: 2,
		CounselorID:   &counselorID,
	})

	year := models.Sophomore
	gpa := 4
	organization := "Robotics"
	interest := "Engineering"
	highSchool := "Brookline High School"
	testScore := "1350"
	major := "Design"
	minor := "Economics"
	geoPrefs := "Northeast"

	updated, err := repo.UpdateStudent(ctx, created.ID, models.UpdateStudentBody{
		Year:                  &year,
		GPA:                   &gpa,
		Organization:          &organization,
		Interest:              &interest,
		HighSchool:            &highSchool,
		StandardizedTestScore: &testScore,
		IntendedMajor:         &major,
		IntendedMinor:         &minor,
		GeographicPreferences: &geoPrefs,
	})
	if err != nil {
		t.Fatalf("UpdateStudent failed: %v", err)
	}

	if updated.Year != year {
		t.Errorf("expected Year %q, got %q", year, updated.Year)
	}
	if updated.GPA != gpa {
		t.Errorf("expected GPA %d, got %d", gpa, updated.GPA)
	}
	if updated.Organization == nil || *updated.Organization != organization {
		t.Errorf("expected Organization %q, got %v", organization, updated.Organization)
	}
	if updated.Interest == nil || *updated.Interest != interest {
		t.Errorf("expected Interest %q, got %v", interest, updated.Interest)
	}
	if updated.HighSchool == nil || *updated.HighSchool != highSchool {
		t.Errorf("expected HighSchool %q, got %v", highSchool, updated.HighSchool)
	}
	if updated.StandardizedTestScore == nil || *updated.StandardizedTestScore != testScore {
		t.Errorf("expected StandardizedTestScore %q, got %v", testScore, updated.StandardizedTestScore)
	}
	if updated.IntendedMajor == nil || *updated.IntendedMajor != major {
		t.Errorf("expected IntendedMajor %q, got %v", major, updated.IntendedMajor)
	}
	if updated.IntendedMinor == nil || *updated.IntendedMinor != minor {
		t.Errorf("expected IntendedMinor %q, got %v", minor, updated.IntendedMinor)
	}
	if updated.GeographicPreferences == nil || *updated.GeographicPreferences != geoPrefs {
		t.Errorf("expected GeographicPreferences %q, got %v", geoPrefs, updated.GeographicPreferences)
	}
	// fields not in the update body are untouched
	if updated.ReviewBalance != created.ReviewBalance {
		t.Errorf("expected ReviewBalance %d, got %d", created.ReviewBalance, updated.ReviewBalance)
	}
	if updated.CounselorID == nil || *updated.CounselorID != counselorID {
		t.Errorf("expected CounselorID %s, got %v", counselorID, updated.CounselorID)
	}
}

// fields left nil keep their current values
func TestUpdateStudent_Partial(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, db := setupTestRepo(t)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	organization := "Debate Club"
	created := createTestStudent(t, db, models.Student{
		UserID:        createTestUser(t, db, "Partly Updated Student"),
		Year:          models.Junior,
		Organization:  &organization,
		GPA:           3,
		ReviewBalance: 0,
		CounselorID:   &counselorID,
	})

	gpa := 4
	updated, err := repo.UpdateStudent(ctx, created.ID, models.UpdateStudentBody{GPA: &gpa})
	if err != nil {
		t.Fatalf("UpdateStudent failed: %v", err)
	}

	if updated.GPA != gpa {
		t.Errorf("expected GPA %d, got %d", gpa, updated.GPA)
	}
	if updated.Year != created.Year {
		t.Errorf("expected Year %q to be unchanged, got %q", created.Year, updated.Year)
	}
	if updated.Organization == nil || *updated.Organization != organization {
		t.Errorf("expected Organization %q to be unchanged, got %v", organization, updated.Organization)
	}
}

func TestUpdateStudent_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	gpa := 4
	_, err := repo.UpdateStudent(ctx, uuid.New(), models.UpdateStudentBody{GPA: &gpa})
	if err == nil {
		t.Fatal("expected error updating non-existent student, got nil")
	}
	httpErr, ok := err.(errs.HTTPError)
	if !ok || httpErr.Code != 404 {
		t.Errorf("expected 404 HTTPError, got %v", err)
	}
}
