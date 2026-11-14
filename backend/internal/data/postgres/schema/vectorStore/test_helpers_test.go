package vectorRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
)

const testDimensions = 1024

func setupTestRepo(t *testing.T) (*VectorRepository, *pgxpool.Pool) {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	db := testutils.SetupTestDB(t)

	config := db.Config()
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	simplePool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		t.Fatalf("failed to create simple-protocol pool: %v", err)
	}
	t.Cleanup(simplePool.Close)

	return NewVectorRepository(simplePool), db
}

// unitVector returns a 1024-dim vector pointing along one axis, so similarities are easy to predict:
// the same axis = similarity 1, different axes = similarity 0.
func unitVector(axis int) []float32 {
	v := make([]float32, testDimensions)
	v[axis] = 1
	return v
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

// createTestMessage inserts a chat message and returns its id.
func createTestMessage(t *testing.T, db *pgxpool.Pool, senderID, recipientID uuid.UUID, message string) uuid.UUID {
	t.Helper()

	var id uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.chat_messages (sender_id, recipient_id, message) VALUES ($1, $2, $3) RETURNING id`,
		senderID, recipientID, message).Scan(&id); err != nil {
		t.Fatalf("setup createTestMessage failed: %v", err)
	}
	return id
}

// createTestTodo inserts a todo item for a student and returns its id.
func createTestTodo(t *testing.T, db *pgxpool.Pool, studentID uuid.UUID, description string) uuid.UUID {
	t.Helper()

	var id uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.todo_items (student_id, user_id, todo_description) VALUES ($1, $1, $2) RETURNING id`,
		studentID, description).Scan(&id); err != nil {
		t.Fatalf("setup createTestTodo failed: %v", err)
	}
	return id
}

// linkVector stores an embedding and links it to a message or todo, like the embed-vectors edge function does.
func linkVector(t *testing.T, db *pgxpool.Pool, embedding []float32, messageID *uuid.UUID, todoID *uuid.UUID) {
	t.Helper()

	var vectorID int64
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.vector_table (embedding) VALUES ($1::extensions.vector) RETURNING vector_id`,
		formatVector(embedding)).Scan(&vectorID); err != nil {
		t.Fatalf("setup linkVector insert vector failed: %v", err)
	}
	if _, err := db.Exec(context.Background(),
		`INSERT INTO public.vector_relationship (vector_id, message_id, todo_id) VALUES ($1, $2, $3)`,
		vectorID, messageID, todoID); err != nil {
		t.Fatalf("setup linkVector insert relationship failed: %v", err)
	}
}
