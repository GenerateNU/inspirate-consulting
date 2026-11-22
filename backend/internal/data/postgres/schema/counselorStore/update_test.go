package counselorRepository

import (
	"context"
	"testing"
	"time"

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

	link := "https://example.com/counselor"
	expiresAt := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	updated, err := repo.UpdateCounselor(ctx, counselorID, models.UpdateCounselorBody{
		CustomLink: models.NewNullable(link),
		ExpiresAt:  models.NewNullable(expiresAt),
	})
	if err != nil {
		t.Fatalf("UpdateCounselor failed: %v", err)
	}

	if updated.ID != counselorID {
		t.Errorf("expected ID %s, got %s", counselorID, updated.ID)
	}
	if updated.UserID == uuid.Nil {
		t.Error("expected UserID to be set")
	}
	if updated.CustomLink == nil || *updated.CustomLink != link {
		t.Errorf("expected custom link %q, got %v", link, updated.CustomLink)
	}
	if updated.ExpiresAt == nil || !updated.ExpiresAt.Equal(expiresAt) {
		t.Errorf("expected expires_at %v, got %v", expiresAt, updated.ExpiresAt)
	}

	// omitted fields keep their stored values
	kept, err := repo.UpdateCounselor(ctx, counselorID, models.UpdateCounselorBody{})
	if err != nil {
		t.Fatalf("UpdateCounselor (omit) failed: %v", err)
	}
	if kept.CustomLink == nil || kept.ExpiresAt == nil {
		t.Errorf("expected omitted fields to be kept, got link=%v expires_at=%v", kept.CustomLink, kept.ExpiresAt)
	}

	// null clears both fields
	cleared, err := repo.UpdateCounselor(ctx, counselorID, models.UpdateCounselorBody{
		CustomLink: models.Null[string](),
		ExpiresAt:  models.Null[time.Time](),
	})
	if err != nil {
		t.Fatalf("UpdateCounselor (clear) failed: %v", err)
	}
	if cleared.CustomLink != nil || cleared.ExpiresAt != nil {
		t.Errorf("expected fields to be cleared, got link=%v expires_at=%v", cleared.CustomLink, cleared.ExpiresAt)
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
