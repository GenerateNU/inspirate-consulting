package auth

import (
	"context"
	"os"

	"github.com/google/uuid"
)

type studentIDContextKey string

const contextKeyStudentID studentIDContextKey = "student_id"

const DefaultTestStudentID = "00000000-0000-0000-0000-000000000002"

// WithStudentID stores the authenticated student's id in the request context.
func WithStudentID(ctx context.Context, studentID string) context.Context {
	if studentID == "" {
		studentID = DefaultTestStudentID
	}
	if _, err := uuid.Parse(studentID); err != nil {
		studentID = DefaultTestStudentID
	}
	return context.WithValue(ctx, contextKeyStudentID, studentID)
}

// GetStudentID reads the authenticated student's id from context.
func GetStudentID(ctx context.Context) string {
	studentID, _ := ctx.Value(contextKeyStudentID).(string)
	if studentID != "" {
		if _, err := uuid.Parse(studentID); err == nil {
			return studentID
		}
	}

	// Fallback: if TEST_MODE is enabled, return a fixed test student id so
	// handlers can operate without real auth during local development.
	// This keeps handler code unchanged when auth is later wired up.
	if os.Getenv("TEST_MODE") == "true" {
		return DefaultTestStudentID
	}

	return ""
}
