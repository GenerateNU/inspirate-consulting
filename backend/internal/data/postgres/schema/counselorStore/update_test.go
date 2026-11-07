package counselorRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func TestUpdateCounselor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewCounselorRepository(db)
	ctx := context.Background()

	counselorID := createTestCounselor(t, db)

	updated, err := repo.UpdateCounselor(ctx, counselorID, models.UpdateCounselorBody{})
	if err != nil {
		t.Fatalf("UpdateCounselor failed: %v", err)
	}

	if updated.ID != counselorID {
		t.Errorf("expected ID %s, got %s", counselorID, updated.ID)
	}
	if updated.UserID == uuid.Nil {
		t.Error("expected UserID to be set")
	}
}

func TestUpdateCounselor_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewCounselorRepository(db)
	ctx := context.Background()

	_, err := repo.UpdateCounselor(ctx, uuid.New(), models.UpdateCounselorBody{})
	if err == nil {
		t.Fatal("expected error updating non-existent counselor, got nil")
	}
	httpErr, ok := err.(errs.HTTPError)
	if !ok || httpErr.Code != 404 {
		t.Errorf("expected 404 HTTPError, got %v", err)
	}
}
