package counselorRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
)

// test list function returns only the counselor's open reviews, oldest first
func TestListPendingReviews(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewCounselorRepository(db)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	otherCounselorID := createTestCounselor(t, db)

	student := createTestStudent(t, db, counselorID, "Alice", "junior")
	otherStudent := createTestStudent(t, db, otherCounselorID, "Bob", "senior")

	essayA := createTestEssay(t, db, student, "personal-statement")
	essayB := createTestEssay(t, db, student, "supplemental")
	otherEssay := createTestEssay(t, db, otherStudent, "personal-statement")

	openA := createTestReview(t, db, student, essayA, false)
	openB := createTestReview(t, db, student, essayB, false)
	// a completed review and another counselor's review should not show up
	createTestReview(t, db, student, essayA, true)
	createTestReview(t, db, otherStudent, otherEssay, false)

	results, err := repo.ListPendingReviews(ctx, counselorID.String())
	if err != nil {
		t.Fatalf("ListPendingReviews failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 pending reviews, got %d", len(results))
	}
	if results[0].TransactionID != openA || results[1].TransactionID != openB {
		t.Errorf("expected reviews %v then %v, got %v then %v", openA, openB, results[0].TransactionID, results[1].TransactionID)
	}

	first := results[0]
	if first.EssayID != essayA {
		t.Errorf("expected EssayID %v, got %v", essayA, first.EssayID)
	}
	if first.EssayType != "personal-statement" {
		t.Errorf("expected EssayType %q, got %q", "personal-statement", first.EssayType)
	}
	if first.EssayLink != "https://example.com/essay" {
		t.Errorf("expected EssayLink %q, got %q", "https://example.com/essay", first.EssayLink)
	}
	if first.StudentID != student {
		t.Errorf("expected StudentID %v, got %v", student, first.StudentID)
	}
	if first.StudentName != "Alice" {
		t.Errorf("expected StudentName %q, got %q", "Alice", first.StudentName)
	}
	if first.RequestedAt.IsZero() {
		t.Error("expected RequestedAt to be set")
	}
}

func TestListPendingReviews_NoReviews(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewCounselorRepository(db)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)
	createTestStudent(t, db, counselorID, "Alice", "junior")

	results, err := repo.ListPendingReviews(ctx, counselorID.String())
	if err != nil {
		t.Fatalf("ListPendingReviews failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 pending reviews, got %d", len(results))
	}
}
