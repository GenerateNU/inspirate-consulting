package notificationPreferencesRepository

import (
	"context"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

func TestUpdateNotificationPreferences(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewNotificationPreferencesRepository(db)
	ctx := context.Background()

	userID := createTestUser(t, db)
	_, err := repo.CreateDefaultNotificationPreferences(ctx, userID)
	if err != nil {
		t.Fatalf("setup CreateDefaultNotificationPreferences failed: %v", err)
	}

	update := models.UpdateNotificationPreferencesRequestBody{
		EmailEnabled:                false,
		WeeklySummaryEnabled:        false,
		DueDateNotificationsEnabled: true,
		DaysBeforeDue:               3,
		NotifyPastDue:               false,
	}

	updated, err := repo.UpdateNotificationPreferences(ctx, userID, update)
	if err != nil {
		t.Fatalf("UpdateNotificationPreferences failed: %v", err)
	}

	if updated.EmailEnabled != false {
		t.Error("expected EmailEnabled to be updated to false")
	}
	if updated.WeeklySummaryEnabled != false {
		t.Error("expected WeeklySummaryEnabled to be updated to false")
	}
	if updated.DaysBeforeDue != 3 {
		t.Errorf("expected DaysBeforeDue to be updated to 3, got %d", updated.DaysBeforeDue)
	}
	if updated.NotifyPastDue != false {
		t.Error("expected NotifyPastDue to be updated to false")
	}

	fetched, err := repo.GetNotificationPreferences(ctx, userID)
	if err != nil {
		t.Fatalf("GetNotificationPreferences after update failed: %v", err)
	}
	if fetched.DaysBeforeDue != 3 {
		t.Errorf("expected persisted DaysBeforeDue of 3, got %d", fetched.DaysBeforeDue)
	}
}

func TestUpdateNotificationPreferences_InvalidDaysBeforeDue(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewNotificationPreferencesRepository(db)
	ctx := context.Background()

	userID := createTestUser(t, db)
	_, err := repo.CreateDefaultNotificationPreferences(ctx, userID)
	if err != nil {
		t.Fatalf("setup CreateDefaultNotificationPreferences failed: %v", err)
	}

	update := models.UpdateNotificationPreferencesRequestBody{
		DaysBeforeDue: -1,
	}

	_, err = repo.UpdateNotificationPreferences(ctx, userID, update)
	if err == nil {
		t.Fatal("expected an error updating to a negative DaysBeforeDue, got nil")
	}
}
