package notificationPreferencesRepository

import (
	"context"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"testing"
)

func TestGetNotificationPreferences_ExistingRow(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewNotificationPreferencesRepository(db)
	ctx := context.Background()

	userID := createTestUser(t, db)
	created, err := repo.CreateDefaultNotificationPreferences(ctx, userID)
	if err != nil {
		t.Fatalf("setup CreateDefaultNotificationPreferences failed: %v", err)
	}

	fetched, err := repo.GetNotificationPreferences(ctx, userID)
	if err != nil {
		t.Fatalf("GetNotificationPreferences failed: %v", err)
	}

	if fetched.UserID != created.UserID {
		t.Errorf("expected UserID %s, got %s", created.UserID, fetched.UserID)
	}
	if fetched.DaysBeforeDue != created.DaysBeforeDue {
		t.Errorf("expected DaysBeforeDue %d, got %d", created.DaysBeforeDue, fetched.DaysBeforeDue)
	}
}

// test Get for a user that doesn't have a preferences row yet, expecting it to create a default row
func TestGetNotificationPreferences_CreatesDefaultOnFirstAccess(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewNotificationPreferencesRepository(db)
	ctx := context.Background()

	userID := createTestUser(t, db)

	fetched, err := repo.GetNotificationPreferences(ctx, userID)
	if err != nil {
		t.Fatalf("expected GetNotificationPreferences to lazily create a default row, got error: %v", err)
	}

	if fetched.UserID != userID {
		t.Errorf("expected UserID %s, got %s", userID, fetched.UserID)
	}
	if !fetched.EmailEnabled {
		t.Error("expected lazily-created row to have EmailEnabled default to true")
	}

	fetchedAgain, err := repo.GetNotificationPreferences(ctx, userID)
	if err != nil {
		t.Fatalf("second GetNotificationPreferences call failed: %v", err)
	}
	if fetchedAgain.UserID != userID {
		t.Errorf("expected UserID %s on second fetch, got %s", userID, fetchedAgain.UserID)
	}
}
