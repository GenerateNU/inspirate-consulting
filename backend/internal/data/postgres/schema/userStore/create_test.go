package userRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// creates a counselor through CreateUser and returns the counselor id.
func createTestCounselor(t *testing.T, repo *UserRepository, db *pgxpool.Pool) uuid.UUID {
	t.Helper()

	input := models.CreateUserInput{}
	input.Body.Name = "Test Counselor"
	input.Body.Email = uuid.NewString() + "@test.com"
	input.Body.Role = models.CounselorRole

	output, err := repo.CreateUser(context.Background(), input, uuid.New())
	if err != nil {
		t.Fatalf("setup createTestCounselor failed: %v", err)
	}

	var counselorID uuid.UUID
	if err := db.QueryRow(context.Background(),
		`SELECT id FROM public.counselor WHERE user_id = $1`, output.Body.User.ID).Scan(&counselorID); err != nil {
		t.Fatalf("setup createTestCounselor lookup failed: %v", err)
	}
	return counselorID
}

func TestCreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	supabaseID := uuid.New()
	key := "pfp-key"

	input := models.CreateUserInput{}
	input.Body.Name = "Aleng123"
	input.Body.Email = "tt@gmail.com"
	input.Body.PfpKey = &key
	input.Body.Role = models.CounselorRole

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	output, err := repo.CreateUser(ctx, input, supabaseID)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if output.Body.User.Name != input.Body.Name {
		t.Errorf("expected name %q, got %q", input.Body.Name, output.Body.User.Name)
	}
	if output.Body.User.SupabaseID != supabaseID {
		t.Errorf("expected supabase_id %q, got %q", supabaseID, output.Body.User.SupabaseID)
	}
	if output.Body.User.PfpKey == nil || *output.Body.User.PfpKey != key {
		t.Errorf("expected pfp_key %q, got %v", key, output.Body.User.PfpKey)
	}
}

func TestCreateUserNullPfpKey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	supabaseID := uuid.New()

	input := models.CreateUserInput{}
	input.Body.Name = "NoPfpUser"
	input.Body.Email = "nopfp@gmail.com"
	input.Body.PfpKey = nil
	input.Body.Role = models.CounselorRole

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	output, err := repo.CreateUser(ctx, input, supabaseID)
	if err != nil {
		t.Fatalf("CreateUser with nil pfp_key failed: %v", err)
	}

	if output.Body.User.PfpKey != nil {
		t.Errorf("expected nil pfp_key, got %q", *output.Body.User.PfpKey)
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	supabaseID := uuid.New()

	input := models.CreateUserInput{}
	input.Body.Name = "DupeUser"
	input.Body.Email = "dupe@gmail.com"
	input.Body.Role = models.CounselorRole

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	_, err := repo.CreateUser(ctx, input, supabaseID)
	if err != nil {
		t.Fatalf("first CreateUser failed: %v", err)
	}

	_, err = repo.CreateUser(ctx, input, supabaseID)
	if err == nil {
		t.Fatal("expected error creating duplicate user with same supabase_id, got nil")
	}
}

// a counselor gets a counselor row linked to the new user
func TestCreateUser_Counselor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	input := models.CreateUserInput{}
	input.Body.Name = "New Counselor"
	input.Body.Email = "newcounselor@gmail.com"
	input.Body.Role = models.CounselorRole

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	output, err := repo.CreateUser(ctx, input, uuid.New())
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	var count int
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM public.counselor WHERE user_id = $1`, output.Body.User.ID).Scan(&count); err != nil {
		t.Fatalf("counselor lookup failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 counselor row, got %d", count)
	}
}

// a student gets a student row with the given fields, assigned to the given counselor
func TestCreateUser_Student(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	counselorID := createTestCounselor(t, repo, db)

	year := models.Junior
	gpa := 4
	organization := "Debate Club"
	interest := "Computer Science"

	input := models.CreateUserInput{}
	input.Body.Name = "New Student"
	input.Body.Email = "newstudent@gmail.com"
	input.Body.Role = models.StudentRole
	input.Body.Year = &year
	input.Body.GPA = &gpa
	input.Body.Organization = &organization
	input.Body.Interest = &interest
	input.Body.CounselorID = &counselorID

	output, err := repo.CreateUser(ctx, input, uuid.New())
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	var student models.Student
	err = db.QueryRow(ctx,
		`SELECT year, gpa, organization, student_interest, review_balance, counselor_id FROM public.student WHERE user_id = $1`,
		output.Body.User.ID,
	).Scan(&student.Year, &student.GPA, &student.Organization, &student.Interest, &student.ReviewBalance, &student.CounselorID)
	if err != nil {
		t.Fatalf("student lookup failed: %v", err)
	}

	if student.Year != year {
		t.Errorf("expected year %q, got %q", year, student.Year)
	}
	if student.GPA != gpa {
		t.Errorf("expected gpa %d, got %d", gpa, student.GPA)
	}
	if student.Organization == nil || *student.Organization != organization {
		t.Errorf("expected organization %q, got %v", organization, student.Organization)
	}
	if student.Interest == nil || *student.Interest != interest {
		t.Errorf("expected student_interest %q, got %v", interest, student.Interest)
	}
	if student.ReviewBalance != 3 {
		t.Errorf("expected review_balance 3, got %d", student.ReviewBalance)
	}
	if student.CounselorID == nil || *student.CounselorID != counselorID {
		t.Errorf("expected counselor_id %s, got %v", counselorID, student.CounselorID)
	}
}

// if the student insert fails, the users row is rolled back too
func TestCreateUser_StudentFailureRollsBackUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	year := models.Senior
	gpa := 3

	input := models.CreateUserInput{}
	input.Body.Name = "Orphan Student"
	input.Body.Email = "orphan@gmail.com"
	input.Body.Role = models.StudentRole
	input.Body.Year = &year
	input.Body.GPA = &gpa

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	supabaseID := uuid.New()

	// a counselor id that doesn't exist fails the student's foreign key
	unknownCounselorID := uuid.New()
	input.Body.CounselorID = &unknownCounselorID

	_, err := repo.CreateUser(ctx, input, supabaseID)
	if err == nil {
		t.Fatal("expected error creating student with unknown counselor, got nil")
	}

	var count int
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM public.users WHERE supabase_id = $1`, supabaseID).Scan(&count); err != nil {
		t.Fatalf("users lookup failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected users row to be rolled back, found %d", count)
	}
}

func TestCreateUser_InvalidRole(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	input := models.CreateUserInput{}
	input.Body.Name = "No Role"
	input.Body.Email = "norole@gmail.com"

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	supabaseID := uuid.New()

	_, err := repo.CreateUser(ctx, input, supabaseID)
	if err == nil {
		t.Fatal("expected error creating user without a role, got nil")
	}

	var count int
	if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM public.users WHERE supabase_id = $1`, supabaseID).Scan(&count); err != nil {
		t.Fatalf("users lookup failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected users row to be rolled back, found %d", count)
	}
}
