package counselorRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// test list function returns only the counselor's students
func TestListStudents(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewCounselorRepository(db)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	otherCounselorID := createTestCounselor(t, db)

	alice := createTestStudent(t, db, counselorID, "Alice", "junior")
	carol := createTestStudent(t, db, counselorID, "Carol", "freshman")
	bob := createTestStudent(t, db, otherCounselorID, "Bob", "senior")

	results, err := repo.ListStudents(ctx, counselorID)
	if err != nil {
		t.Fatalf("ListStudents failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 students, got %d", len(results))
	}

	found := map[uuid.UUID]models.CounselorStudentRow{}
	for _, r := range results {
		found[r.StudentID] = r
	}
	if _, ok := found[bob]; ok {
		t.Error("expected the other counselor's student to NOT be present in results")
	}
	if _, ok := found[carol]; !ok {
		t.Error("expected Carol to be present in results")
	}

	got, ok := found[alice]
	if !ok {
		t.Fatal("expected Alice to be present in results")
	}
	if got.Name != "Alice" {
		t.Errorf("expected Name %q, got %q", "Alice", got.Name)
	}
	if got.Year != models.Junior {
		t.Errorf("expected Year %q, got %q", models.Junior, got.Year)
	}
}

func TestListStudents_NoStudents(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewCounselorRepository(db)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)

	results, err := repo.ListStudents(ctx, counselorID)
	if err != nil {
		t.Fatalf("ListStudents failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 students, got %d", len(results))
	}
}

// test list function returns the per-student pending review, open task, and unread message counts
func TestListStudents_Counts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewCounselorRepository(db)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	studentID := createTestStudent(t, db, counselorID, "Alice", "junior")

	var counselorUserID, studentUserID uuid.UUID
	if err := db.QueryRow(ctx,
		`SELECT user_id FROM public.counselor WHERE id = $1`, counselorID).Scan(&counselorUserID); err != nil {
		t.Fatalf("failed to look up counselor user id: %v", err)
	}
	if err := db.QueryRow(ctx,
		`SELECT user_id FROM public.student WHERE id = $1`, studentID).Scan(&studentUserID); err != nil {
		t.Fatalf("failed to look up student user id: %v", err)
	}

	// Two todos: one open, one completed
	now := time.Now()
	createTestTodoItem(t, db, studentID, nil, nil)
	createTestTodoItem(t, db, studentID, nil, &now)

	// Two reviews: one open, one completed
	essayID := createTestEssay(t, db, studentID, "personal-statement")
	createTestReview(t, db, studentID, essayID, false)
	createTestReview(t, db, studentID, essayID, true)

	// Messages from the student to the counselor: one unread, one already read.
	// A message from the counselor to the student shouldn't count either.
	if _, err := db.Exec(ctx,
		`INSERT INTO public.chat_messages (sender_id, recipient_id, message, read_at)
		 VALUES ($1, $2, 'Hi, can you look at my essay?', NULL),
		        ($1, $2, 'Thanks!', now()),
		        ($2, $1, 'Sure, sending notes soon', NULL)`,
		studentUserID, counselorUserID); err != nil {
		t.Fatalf("failed to insert chat messages: %v", err)
	}

	// Two activities: the newer one should be returned
	if _, err := db.Exec(ctx,
		`INSERT INTO public.activity (student_id, action, created_at)
		 VALUES ($1, 'added_school', now() - interval '1 day'),
		        ($1, 'added_extracurricular', now())`,
		studentID); err != nil {
		t.Fatalf("failed to insert activities: %v", err)
	}

	students, err := repo.ListStudents(ctx, counselorID)
	if err != nil {
		t.Fatalf("ListStudents failed: %v", err)
	}
	if len(students) != 1 {
		t.Fatalf("expected 1 student, got %d", len(students))
	}

	got := students[0]
	if got.OpenTaskCount != 1 {
		t.Errorf("expected OpenTaskCount 1, got %d", got.OpenTaskCount)
	}
	if got.UnreadMessageCount != 1 {
		t.Errorf("expected UnreadMessageCount 1, got %d", got.UnreadMessageCount)
	}
	if got.PendingReviewCount != 1 {
		t.Errorf("expected PendingReviewCount 1, got %d", got.PendingReviewCount)
	}
	if got.LatestActivityAction == nil {
		t.Fatal("expected LatestActivityAction to be set, got nil")
	}
	if *got.LatestActivityAction != models.ActionAddedExtracurricular {
		t.Errorf("expected LatestActivityAction %q, got %q", models.ActionAddedExtracurricular, *got.LatestActivityAction)
	}
	if got.LatestActivityAt == nil {
		t.Error("expected LatestActivityAt to be set, got nil")
	}
}

// test a student with no activity is still listed, with nil latest-activity fields
func TestListStudents_NoActivity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewCounselorRepository(db)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	createTestStudent(t, db, counselorID, "Alice", "junior")

	students, err := repo.ListStudents(ctx, counselorID)
	if err != nil {
		t.Fatalf("ListStudents failed: %v", err)
	}
	if len(students) != 1 {
		t.Fatalf("expected 1 student, got %d", len(students))
	}
	if students[0].LatestActivityAction != nil {
		t.Errorf("expected LatestActivityAction nil, got %q", *students[0].LatestActivityAction)
	}
	if students[0].LatestActivityAt != nil {
		t.Errorf("expected LatestActivityAt nil, got %v", *students[0].LatestActivityAt)
	}
}
