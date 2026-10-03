package notificationPreferencesRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)


func createTestUser(t *testing.T, db *pgxpool.Pool) uuid.UUID {
	t.Helper()

	var userID uuid.UUID
	err := db.QueryRow(
		context.Background(),
		`INSERT INTO public.users (name) VALUES ($1) RETURNING id`,
		"Test User",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	return userID
}