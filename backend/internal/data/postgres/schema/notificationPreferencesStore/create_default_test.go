package notificationPreferencesRepository

import (
	"context"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"testing"
)

func TestCreateDefaultNotificationPreferences(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewNotificationPreferencesRepository(db)
	ctx := context.Background()

	userID := createTestUser(t, db)

	created, err := repo.CreateDefaultNotificationPreferences(ctx, userID)
	if err != nil {
		t.Fatalf("CreateDefaultNotificationPreferences failed: %v", err)
	}

	if created.UserID != userID {
		t.Errorf("expected UserID %s, got %s", userID, created.UserID)
	}
	if !created.EmailEnabled {
		t.Error("expected EmailEnabled to default to true")
	}
	if !created.WeeklySummaryEnabled {
		t.Error("expected WeeklySummaryEnabled to default to true")
	}
	if !created.DueDateNotificationsEnabled {
		t.Error("expected DueDateNotificationsEnabled to default to true")
	}
	if created.DaysBeforeDue != 1 {
		t.Errorf("expected DaysBeforeDue to default to 1, got %d", created.DaysBeforeDue)
	}
	if !created.NotifyPastDue {
		t.Error("expected NotifyPastDue to default to true")
	}
}

func TestCreateDefaultNotificationPreferences_DuplicateUserFails(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewNotificationPreferencesRepository(db)
	ctx := context.Background()

	userID := createTestUser(t, db)

	_, err := repo.CreateDefaultNotificationPreferences(ctx, userID)
	if err != nil {
		t.Fatalf("first CreateDefaultNotificationPreferences failed: %v", err)
	}

	_, err = repo.CreateDefaultNotificationPreferences(ctx, userID)
	if err == nil {
		t.Fatal("expected an error creating a second preferences row for the same user, got nil")
	}
}