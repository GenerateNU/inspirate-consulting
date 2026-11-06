package chatMessageRepository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// setupTestRepo creates a fresh, isolated test database and returns a repository connected to it.
func setupTestRepo(t *testing.T) (*ChatMessageRepository, *pgxpool.Pool) {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	db := testutils.SetupTestDB(t)
	return NewChatMessageRepository(db), db
}

// createTestUser inserts a user directly and returns its id.
func createTestUser(t *testing.T, db *pgxpool.Pool, name string) uuid.UUID {
	t.Helper()

	var userID uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.users (name, email) VALUES ($1, $2) RETURNING id`,
		name, uuid.NewString()+"@test.com").Scan(&userID); err != nil {
		t.Fatalf("setup createTestUser(%q) failed: %v", name, err)
	}
	return userID
}

// insertMessage inserts a message with an explicit created_at so tests control ordering.
func insertMessage(t *testing.T, db *pgxpool.Pool, senderID, recipientID uuid.UUID, message string, createdAt time.Time) models.ChatMessage {
	t.Helper()

	var m models.ChatMessage
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.chat_messages (sender_id, recipient_id, message, created_at)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, sender_id, recipient_id, message, created_at, edited_at, read_at`,
		senderID, recipientID, message, createdAt,
	).Scan(&m.ID, &m.SenderID, &m.RecipientID, &m.Message, &m.CreatedAt, &m.EditedAt, &m.ReadAt); err != nil {
		t.Fatalf("setup insertMessage failed: %v", err)
	}
	return m
}

// getReadAt reads a message's read_at straight from the table.
func getReadAt(t *testing.T, db *pgxpool.Pool, id uuid.UUID) *time.Time {
	t.Helper()

	var readAt *time.Time
	if err := db.QueryRow(context.Background(),
		`SELECT read_at FROM public.chat_messages WHERE id = $1`, id).Scan(&readAt); err != nil {
		t.Fatalf("getReadAt failed: %v", err)
	}
	return readAt
}

// baseTime is a fixed reference point for message timestamps.
var baseTime = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
