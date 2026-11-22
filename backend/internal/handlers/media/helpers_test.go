package media

import (
	"context"

	"inspirate-consulting/internal/auth"

	"github.com/google/uuid"
)

func ptr[T any](v T) *T { return &v }

func studentCtx(studentID uuid.UUID) context.Context {
	return context.WithValue(auth.WithStudentID(context.Background(), studentID), "Role", "student") //nolint:staticcheck // matches how auth middleware stores the role
}

func counselorCtx(counselorID uuid.UUID) context.Context {
	return context.WithValue(auth.WithCounselorID(context.Background(), counselorID), "Role", "counselor") //nolint:staticcheck // matches how auth middleware stores the role
}
