package personalCollegeApplicationRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

func addTestCycle(t *testing.T, db *pgxpool.Pool, collegeID int64, cycle string) {
	t.Helper()
	if _, err := db.Exec(context.Background(),
		`INSERT INTO public.college_application_cycles (global_college_id, cycle, deadline) VALUES ($1, $2, now() + interval '30 days')`,
		collegeID, cycle); err != nil {
		t.Fatalf("failed to add test cycle: %v", err)
	}
}

// collegeID nil = a general / Common App essay
func addTestEssay(t *testing.T, db *pgxpool.Pool, studentID uuid.UUID, collegeID *int64, status string) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.essays (student_id, type, college_id, link_to_content, status)
		 VALUES ($1, 'supplement', $2, 'https://example.com/essay', $3) RETURNING id`,
		studentID, collegeID, status).Scan(&id); err != nil {
		t.Fatalf("failed to add test essay: %v", err)
	}
	return id
}

// creates a real student (the review ledger has an FK to student) and returns its id
func createTestStudent(t *testing.T, db *pgxpool.Pool) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	var userID, counselorID, studentID uuid.UUID
	if err := db.QueryRow(ctx,
		`INSERT INTO public.users (name, email) VALUES ('college list test', $1) RETURNING id`, uuid.NewString()+"@test.com").Scan(&userID); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	if err := db.QueryRow(ctx, `INSERT INTO public.counselor (user_id) VALUES ($1) RETURNING id`, userID).Scan(&counselorID); err != nil {
		t.Fatalf("failed to create test counselor: %v", err)
	}
	if err := db.QueryRow(ctx,
		`INSERT INTO public.student (user_id, year, gpa, review_balance, counselor_id) VALUES ($1, 'senior', 4, 0, $2) RETURNING id`,
		userID, counselorID).Scan(&studentID); err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}
	return studentID
}

func essayExists(t *testing.T, db *pgxpool.Pool, id uuid.UUID) (bool, string) {
	t.Helper()
	var status string
	err := db.QueryRow(context.Background(), `SELECT status::text FROM public.essays WHERE id = $1`, id).Scan(&status)
	if err != nil {
		return false, ""
	}
	return true, status
}

func TestApplication_ComputedFields(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()
	studentID := uuid.New()

	collegeID := createTestGlobalCollege(t, db, "Computed University", "Computed City, CC")
	addTestCycle(t, db, collegeID, models.ApplicationCycleED)

	created, err := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID, ApplicationType: models.ApplicationCycleED, Category: "reach",
	})
	if err != nil {
		t.Fatalf("CreatePersonalCollegeApplication failed: %v", err)
	}
	if !created.CycleOffered {
		t.Error("expected cycle_offered=true for an offered cycle")
	}
	if created.EssayProgress != models.EssayProgressNotStarted {
		t.Errorf("expected essay_progress not_started with no essays, got %q", created.EssayProgress)
	}

	progress := func() string {
		t.Helper()
		apps, err := repo.ListPersonalCollegeApplicationsByStudentID(ctx, studentID)
		if err != nil || len(apps) != 1 {
			t.Fatalf("list failed: %v (%d apps)", err, len(apps))
		}
		return apps[0].EssayProgress
	}

	// general essays and archived essays don't count
	addTestEssay(t, db, studentID, nil, "Draft")
	addTestEssay(t, db, studentID, &collegeID, "Archived")
	if got := progress(); got != models.EssayProgressNotStarted {
		t.Errorf("expected not_started ignoring general/archived essays, got %q", got)
	}

	draft := addTestEssay(t, db, studentID, &collegeID, "Draft")
	addTestEssay(t, db, studentID, &collegeID, "Submitted")
	if got := progress(); got != models.EssayProgressInProgress {
		t.Errorf("expected in_progress with one draft + one submitted, got %q", got)
	}

	if _, err := db.Exec(ctx, `UPDATE public.essays SET status = 'Submitted' WHERE id = $1`, draft); err != nil {
		t.Fatal(err)
	}
	if got := progress(); got != models.EssayProgressCompleted {
		t.Errorf("expected completed with all essays submitted, got %q", got)
	}

	// counselor removes ED: the application is kept but flagged
	if _, err := db.Exec(ctx, `DELETE FROM public.college_application_cycles WHERE global_college_id = $1`, collegeID); err != nil {
		t.Fatal(err)
	}
	apps, err := repo.ListPersonalCollegeApplicationsByStudentID(ctx, studentID)
	if err != nil || len(apps) != 1 {
		t.Fatalf("list failed: %v", err)
	}
	if apps[0].CycleOffered {
		t.Error("expected cycle_offered=false after the college stopped offering ED")
	}
}

func TestUpdatePersonalCollegeApplication_InterestAndStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()
	studentID := uuid.New()

	collegeID := createTestGlobalCollege(t, db, "Interest University", "Interest City, IC")
	created, err := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID, ApplicationType: models.ApplicationCycleRD, Category: "target",
	})
	if err != nil {
		t.Fatalf("CreatePersonalCollegeApplication failed: %v", err)
	}

	interest := models.InterestTopChoice
	updated, err := repo.UpdatePersonalCollegeApplication(ctx, studentID, created.ID, models.UpdatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID:   collegeID,
		ApplicationType:   models.ApplicationCycleED2,
		Category:          "reach",
		Interest:          &interest,
		ApplicationStatus: models.ApplicationStatusSubmitted,
	})
	if err != nil {
		t.Fatalf("UpdatePersonalCollegeApplication failed: %v", err)
	}
	if updated.Interest == nil || *updated.Interest != interest {
		t.Errorf("expected interest %q, got %v", interest, updated.Interest)
	}
	if updated.ApplicationStatus != models.ApplicationStatusSubmitted || updated.ApplicationType != models.ApplicationCycleED2 {
		t.Errorf("expected ED2/submitted, got %s/%s", updated.ApplicationType, updated.ApplicationStatus)
	}

	// omitting interest clears it
	cleared, err := repo.UpdatePersonalCollegeApplication(ctx, studentID, created.ID, models.UpdatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID:   collegeID,
		ApplicationType:   models.ApplicationCycleED2,
		Category:          "reach",
		ApplicationStatus: models.ApplicationStatusInProgress,
	})
	if err != nil {
		t.Fatalf("UpdatePersonalCollegeApplication (clear) failed: %v", err)
	}
	if cleared.Interest != nil {
		t.Errorf("expected interest cleared, got %q", *cleared.Interest)
	}

	_, err = repo.UpdatePersonalCollegeApplication(ctx, studentID, created.ID, models.UpdatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID, ApplicationType: "RD", Category: "reach", ApplicationStatus: "bogus",
	})
	if err == nil {
		t.Error("expected an error for an invalid application_status")
	}
}

func TestDeletePersonalCollegeApplication_RemovesCollegeEssays(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()
	studentID := createTestStudent(t, db)
	otherStudentID := createTestStudent(t, db)

	collegeID := createTestGlobalCollege(t, db, "Essay Delete University", "Essay City, EC")
	otherCollegeID := createTestGlobalCollege(t, db, "Kept College", "Kept City, KC")

	app, err := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID, ApplicationType: "RD", Category: "reach",
	})
	if err != nil {
		t.Fatalf("CreatePersonalCollegeApplication failed: %v", err)
	}

	collegeEssay := addTestEssay(t, db, studentID, &collegeID, "Draft")
	reviewedEssay := addTestEssay(t, db, studentID, &collegeID, "Review")
	generalEssay := addTestEssay(t, db, studentID, nil, "Draft")
	otherCollegeEssay := addTestEssay(t, db, studentID, &otherCollegeID, "Draft")
	otherStudentEssay := addTestEssay(t, db, otherStudentID, &collegeID, "Draft")

	if _, err := db.Exec(ctx,
		`INSERT INTO public.essay_review_transaction (subtotal, student_id, entry_type, essay_id) VALUES (-1, $1, 'spend', $2)`,
		studentID, reviewedEssay); err != nil {
		t.Fatalf("failed to add review transaction: %v", err)
	}

	if err := repo.DeletePersonalCollegeApplication(ctx, studentID, app.ID); err != nil {
		t.Fatalf("DeletePersonalCollegeApplication failed: %v", err)
	}

	if ok, _ := essayExists(t, db, collegeEssay); ok {
		t.Error("expected the student's essay for the removed college to be deleted")
	}
	if ok, status := essayExists(t, db, reviewedEssay); !ok || status != "Archived" {
		t.Errorf("expected the reviewed essay to be kept and archived, got exists=%v status=%q", ok, status)
	}
	for name, id := range map[string]uuid.UUID{"general": generalEssay, "other college": otherCollegeEssay, "other student": otherStudentEssay} {
		if ok, status := essayExists(t, db, id); !ok || status != "Draft" {
			t.Errorf("expected %s essay untouched, got exists=%v status=%q", name, ok, status)
		}
	}
}

func TestDeletePersonalCollegeApplication_KeepsEssaysWhenCollegeStillListed(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()
	studentID := uuid.New()

	collegeID := createTestGlobalCollege(t, db, "Twice Listed University", "Twice City, TC")
	input := models.CreatePersonalCollegeApplicationRequestBody{GlobalCollegeID: collegeID, ApplicationType: "RD", Category: "reach"}
	first, err := repo.CreatePersonalCollegeApplication(ctx, studentID, input)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreatePersonalCollegeApplication(ctx, studentID, input); err != nil {
		t.Fatal(err)
	}
	essay := addTestEssay(t, db, studentID, &collegeID, "Draft")

	if err := repo.DeletePersonalCollegeApplication(ctx, studentID, first.ID); err != nil {
		t.Fatalf("DeletePersonalCollegeApplication failed: %v", err)
	}
	if ok, _ := essayExists(t, db, essay); !ok {
		t.Error("expected essay kept while another application to the same college remains")
	}
}

func TestDeletePersonalCollegeApplication_ClosesRankGap(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()
	studentID := uuid.New()

	ids := make([]int64, 3)
	for i, name := range []string{"Rank One U", "Rank Two U", "Rank Three U"} {
		collegeID := createTestGlobalCollege(t, db, name, "Rank City, RC")
		app, err := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: collegeID, ApplicationType: "RD", Category: "target",
		})
		if err != nil {
			t.Fatal(err)
		}
		rank := i + 1
		if _, err := repo.UpdateApplicationRank(ctx, studentID, app.ID, &rank); err != nil {
			t.Fatal(err)
		}
		ids[i] = app.ID
	}

	if err := repo.DeletePersonalCollegeApplication(ctx, studentID, ids[0]); err != nil {
		t.Fatalf("DeletePersonalCollegeApplication failed: %v", err)
	}

	apps, err := repo.ListPersonalCollegeApplicationsByStudentID(ctx, studentID)
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 2 {
		t.Fatalf("expected 2 applications left, got %d", len(apps))
	}
	for i, app := range apps {
		if app.Rank == nil || *app.Rank != i+1 {
			t.Errorf("expected rank %d for application %d, got %v", i+1, app.ID, app.Rank)
		}
	}
	if apps[0].ID != ids[1] {
		t.Errorf("expected former #2 to be #1 now, got application %d", apps[0].ID)
	}
}
