package userRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/errs"
)

func createSessionTestUser(t *testing.T, db *pgxpool.Pool, name string, supabaseID uuid.UUID) uuid.UUID {
	t.Helper()

	var userID uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.users (name, supabase_id, email) VALUES ($1, $2, $3) RETURNING id`, name, supabaseID, uuid.NewString()+"@test.com").Scan(&userID); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return userID
}

func TestFetchSessionUserBySupabaseID_Counselor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	supabaseID := uuid.New()
	userID := createSessionTestUser(t, db, "Session Counselor", supabaseID)
	var counselorID uuid.UUID
	if err := db.QueryRow(ctx,
		`INSERT INTO public.counselor (user_id) VALUES ($1) RETURNING id`, userID).Scan(&counselorID); err != nil {
		t.Fatalf("failed to create test counselor: %v", err)
	}

	user, err := repo.FetchSessionUserBySupabaseID(ctx, supabaseID)
	if err != nil {
		t.Fatalf("FetchSessionUserBySupabaseID failed: %v", err)
	}

	if user.ID != userID {
		t.Errorf("expected id %v, got %v", userID, user.ID)
	}
	if user.CounselorID == nil || *user.CounselorID != counselorID {
		t.Errorf("expected counselor_id %v, got %v", counselorID, user.CounselorID)
	}
	if user.StudentID != nil {
		t.Errorf("expected no student_id for a counselor, got %v", *user.StudentID)
	}
}

func TestFetchSessionUserBySupabaseID_Student(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	counselorUserID := createSessionTestUser(t, db, "Assigned Counselor", uuid.New())
	var counselorID uuid.UUID
	if err := db.QueryRow(ctx,
		`INSERT INTO public.counselor (user_id) VALUES ($1) RETURNING id`, counselorUserID).Scan(&counselorID); err != nil {
		t.Fatalf("failed to create test counselor: %v", err)
	}

	supabaseID := uuid.New()
	userID := createSessionTestUser(t, db, "Session Student", supabaseID)
	var studentID uuid.UUID
	if err := db.QueryRow(ctx,
		`INSERT INTO public.student (user_id, year, gpa, review_balance, counselor_id)
		 VALUES ($1, 'junior', 4, 0, $2) RETURNING id`, userID, counselorID).Scan(&studentID); err != nil {
		t.Fatalf("failed to create test student: %v", err)
	}

	user, err := repo.FetchSessionUserBySupabaseID(ctx, supabaseID)
	if err != nil {
		t.Fatalf("FetchSessionUserBySupabaseID failed: %v", err)
	}

	if user.StudentID == nil || *user.StudentID != studentID {
		t.Errorf("expected student_id %v, got %v", studentID, user.StudentID)
	}
	if user.CounselorID != nil {
		t.Errorf("expected no counselor_id for a student, got %v", *user.CounselorID)
	}
}

func TestFetchSessionUserBySupabaseID_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.FetchSessionUserBySupabaseID(context.Background(), uuid.New())
	httpErr, ok := err.(errs.HTTPError)
	if !ok || httpErr.Code != 404 {
		t.Fatalf("expected a 404 error, got %v", err)
	}
}
