package essays

import (
	"context"

	"inspirate-consulting/internal/auth"

	"github.com/google/uuid"
)

// studentRoleCtx returns a context for a logged-in student, with the role set
// so a real Authorizer enforces ownership.
func studentRoleCtx(studentID uuid.UUID) context.Context {
	return context.WithValue(auth.WithStudentID(context.Background(), studentID), "Role", "student") //nolint:staticcheck // matches how auth middleware stores the role
}

// realAuthorizer returns an Authorizer that enforces the student ownership
// branch. It has no database, so it must not be used for counselor checks.
func realAuthorizer() *auth.Authorizer {
	return auth.NewAuthorizer(nil)
}
